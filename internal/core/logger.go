package core

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// LogLevel represents the logging level
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// LogFormat represents the logging format
type LogFormat string

const (
	LogFormatText   LogFormat = "text"
	LogFormatJSON   LogFormat = "json"
	LogFormatCustom LogFormat = "custom"
)

type LoggerConfig struct {
	Level      LogLevel
	Format     LogFormat
	AddSource  bool
	Output     io.Writer
	AppName    string
	AppVersion string
}

type CustomHandler struct {
	opts   *slog.HandlerOptions
	output io.Writer
	attrs  []slog.Attr
	groups []string
}

func NewCustomHandler(output io.Writer, opts *slog.HandlerOptions) *CustomHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &CustomHandler{
		opts:   opts,
		output: output,
		attrs:  []slog.Attr{},
		groups: []string{},
	}
}

// Enabled reports whether the handler handles records at the given level
func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	timestamp := r.Time.Format("02/01/2006 15:04:05")
	level := strings.ToUpper(r.Level.String())

	firstLine := fmt.Sprintf("%s %s:", timestamp, level)

	var secondLine string

	if h.opts.AddSource && r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()

		filename := filepath.Base(f.File)
		source := fmt.Sprintf("%s:%d", filename, f.Line)

		secondLine = fmt.Sprintf("    [%s] %s", source, r.Message)
	} else {
		secondLine = fmt.Sprintf("    %s", r.Message)
	}

	var attrs []string
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, fmt.Sprintf("%s=%v", a.Key, a.Value))
		return true
	})

	// Add handler attributes
	for _, attr := range h.attrs {
		attrs = append(attrs, fmt.Sprintf("%s=%v", attr.Key, attr.Value))
	}

	if len(attrs) > 0 {
		secondLine += " " + strings.Join(attrs, " ")
	}

	// Write both lines
	output := fmt.Sprintf("%s\n%s\n", firstLine, secondLine)
	_, err := h.output.Write([]byte(output))
	return err
}

// WithAttrs returns a new handler with the given attributes
func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)

	return &CustomHandler{
		opts:   h.opts,
		output: h.output,
		attrs:  newAttrs,
		groups: h.groups,
	}
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	newGroups := make([]string, len(h.groups)+1)
	copy(newGroups, h.groups)
	newGroups[len(h.groups)] = name

	return &CustomHandler{
		opts:   h.opts,
		output: h.output,
		attrs:  h.attrs,
		groups: newGroups,
	}
}

func NewLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level:      getLogLevelFromEnv(),
		Format:     getLogFormatFromEnv(),
		AddSource:  getAddSourceFromEnv(),
		Output:     os.Stdout,
		AppName:    getAppNameFromEnv(),
		AppVersion: getAppVersionFromEnv(),
	}
}

func NewLogger(config *LoggerConfig) *slog.Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level:     getSlogLevel(config.Level),
		AddSource: config.AddSource,
	}

	switch config.Format {
	case LogFormatJSON:
		handler = slog.NewJSONHandler(config.Output, opts)
	case LogFormatCustom:
		handler = NewCustomHandler(config.Output, opts)
	default:
		handler = NewCustomHandler(config.Output, opts)
	}

	return slog.New(handler)
}

func NewDefaultLogger() *slog.Logger {
	config := NewLoggerConfig()
	return NewLogger(config)
}

// Helper functions to get configuration from environment variables

func getLogLevelFromEnv() LogLevel {
	level := strings.ToLower(os.Getenv("LOG_LEVEL"))
	switch level {
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "warn", "warning":
		return LogLevelWarn
	case "error":
		return LogLevelError
	default:
		// Default to info in production, debug in development
		if isDevelopment() {
			return LogLevelDebug
		}
		return LogLevelInfo
	}
}

func getLogFormatFromEnv() LogFormat {
	format := strings.ToLower(os.Getenv("LOG_FORMAT"))
	switch format {
	case "json":
		return LogFormatJSON
	case "text":
		return LogFormatText
	case "custom":
		return LogFormatCustom
	default:
		// Default to custom format for better readability
		return LogFormatCustom
	}
}

func getAddSourceFromEnv() bool {
	addSource := strings.ToLower(os.Getenv("LOG_ADD_SOURCE"))
	switch addSource {
	case "true", "1", "yes":
		return true
	case "false", "0", "no":
		return false
	default:
		// Default to true in development, false in production
		return isDevelopment()
	}
}

func getAppNameFromEnv() string {
	if name := os.Getenv("APP_NAME"); name != "" {
		return name
	}
	return "pandoragym-api"
}

func getAppVersionFromEnv() string {
	if version := os.Getenv("APP_VERSION"); version != "" {
		return version
	}
	return "1.0.0"
}

func isDevelopment() bool {
	env := strings.ToLower(os.Getenv("GO_ENV"))
	return env == "development" || env == "dev" || env == ""
}

func getSlogLevel(level LogLevel) slog.Level {
	switch level {
	case LogLevelDebug:
		return slog.LevelDebug
	case LogLevelInfo:
		return slog.LevelInfo
	case LogLevelWarn:
		return slog.LevelWarn
	case LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// LoggerMiddleware creates a logger with request context
func LoggerWithContext(logger *slog.Logger, requestID string) *slog.Logger {
	return logger.With(
		slog.String("request_id", requestID),
	)
}

// LoggerWithService creates a logger with service context
func LoggerWithService(logger *slog.Logger, service string) *slog.Logger {
	return logger.With(
		slog.String("service", service),
	)
}

// LoggerWithUser creates a logger with user context
func LoggerWithUser(logger *slog.Logger, userID string) *slog.Logger {
	return logger.With(
		slog.String("user_id", userID),
	)
}
