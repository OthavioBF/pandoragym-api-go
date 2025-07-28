# Logging Configuration

This document describes the improved logging system for the PandoraGym API.

## Overview

The logging system uses Go's structured logging (`log/slog`) with environment-based configuration and contextual logging capabilities.

## Configuration

### Environment Variables

| Variable | Description | Default | Options |
|----------|-------------|---------|---------|
| `LOG_LEVEL` | Logging level | `debug` (dev), `info` (prod) | `debug`, `info`, `warn`, `error` |
| `LOG_FORMAT` | Output format | `custom` (dev), `json` (prod) | `text`, `json`, `custom` |
| `LOG_ADD_SOURCE` | Include source file info | `true` (dev), `false` (prod) | `true`, `false` |
| `APP_NAME` | Application name in logs | `pandoragym-api` | Any string |
| `APP_VERSION` | Application version in logs | `1.0.0` | Any string |
| `GO_ENV` | Environment mode | `development` | `development`, `production` |

### Example .env Configuration

```env
# Development
LOG_LEVEL=debug
LOG_FORMAT=custom
LOG_ADD_SOURCE=true
APP_NAME=pandoragym-api
APP_VERSION=1.0.0
GO_ENV=development

# Production
LOG_LEVEL=info
LOG_FORMAT=json
LOG_ADD_SOURCE=false
APP_NAME=pandoragym-api
APP_VERSION=1.0.0
GO_ENV=production
```

## Usage Examples

### Basic Logging

```go
logger := core.NewDefaultLogger()

logger.Info("Application started")
logger.Error("Something went wrong", slog.String("error", err.Error()))
```

### Contextual Logging

```go
// Add request context
requestLogger := core.LoggerWithContext(logger, requestID)
requestLogger.Info("Processing request")

// Add service context
serviceLogger := core.LoggerWithService(logger, "user-service")
serviceLogger.Debug("Executing business logic")

// Add user context
userLogger := core.LoggerWithUser(logger, userID)
userLogger.Info("User action completed")
```

### Structured Logging

```go
logger.Info("User created",
    slog.String("user_id", user.ID),
    slog.String("email", user.Email),
    slog.String("role", user.Role),
    slog.Time("created_at", user.CreatedAt),
)
```

### Grouped Attributes

```go
logger.Info("Database operation completed",
    slog.Group("database",
        slog.String("table", "users"),
        slog.String("operation", "INSERT"),
        slog.Int("rows_affected", 1),
    ),
    slog.Group("performance",
        slog.Duration("duration", duration),
        slog.Int("connections", activeConnections),
    ),
)
```

### Error Logging with Context

```go
logger.Error("Failed to create user",
    slog.String("error", err.Error()),
    slog.String("error_type", "validation_error"),
    slog.Group("request",
        slog.String("method", "POST"),
        slog.String("path", "/auth/register/student"),
        slog.String("user_agent", userAgent),
    ),
    slog.Group("validation",
        slog.String("field", "email"),
        slog.String("value", email),
        slog.String("constraint", "unique"),
    ),
)
```

## Log Levels

### Debug
- Detailed information for debugging
- Only shown in development
- Use for: Variable values, function entry/exit, detailed flow

```go
logger.Debug("Processing user registration",
    slog.String("email", email),
    slog.String("role", role),
)
```

### Info
- General information about application flow
- Shown in all environments
- Use for: Successful operations, important state changes

```go
logger.Info("User registered successfully",
    slog.String("user_id", userID),
    slog.String("role", role),
)
```

### Warn
- Warning conditions that don't stop execution
- Use for: Deprecated features, fallback scenarios, recoverable errors

```go
logger.Warn("JWT_SECRET not set, using default value",
    slog.Bool("production", isProduction),
)
```

### Error
- Error conditions that affect functionality
- Use for: Failed operations, exceptions, critical issues

```go
logger.Error("Failed to connect to database",
    slog.String("error", err.Error()),
    slog.String("database_url", dbURL),
)
```

## Output Formats

### Custom Format (Development - Two Lines)
```
2025-07-28 19:00:00 INFO:
    [auth_handlers.go:65] User registered successfully user_id=123 role=student app.name=pandoragym-api app.version=1.0.0
```

### Text Format (Standard slog)
```
time=2025-07-28T19:00:00.000Z level=INFO source=internal/api/auth_handlers.go:65 msg="User registered successfully" app.name=pandoragym-api app.version=1.0.0 user_id=123 role=student
```

### JSON Format (Production)
```json
{
  "time": "2025-07-28T19:00:00.000Z",
  "level": "INFO",
  "source": {
    "function": "github.com/othavioBF/pandoragym-go-api/internal/api.(*API).RegisterStudent",
    "file": "internal/api/auth_handlers.go",
    "line": 65
  },
  "msg": "User registered successfully",
  "app": {
    "name": "pandoragym-api",
    "version": "1.0.0"
  },
  "user_id": "123",
  "role": "student"
}
```

## Best Practices

### 1. Use Appropriate Log Levels
- **Debug**: Detailed debugging information
- **Info**: Normal application flow
- **Warn**: Unusual but handled situations
- **Error**: Error conditions

### 2. Include Relevant Context
```go
// Good
logger.Info("User login successful",
    slog.String("user_id", userID),
    slog.String("ip_address", clientIP),
    slog.Duration("login_duration", duration),
)

// Avoid
logger.Info("Login successful")
```

### 3. Use Structured Attributes
```go
// Good
logger.Error("Database query failed",
    slog.String("error", err.Error()),
    slog.String("query", "SELECT * FROM users"),
    slog.Int("timeout_seconds", 30),
)

// Avoid
logger.Error(fmt.Sprintf("Database query failed: %v", err))
```

### 4. Group Related Attributes
```go
logger.Info("HTTP request processed",
    slog.Group("request",
        slog.String("method", method),
        slog.String("path", path),
        slog.Int("status", status),
    ),
    slog.Group("performance",
        slog.Duration("duration", duration),
        slog.Int("response_size", size),
    ),
)
```

### 5. Don't Log Sensitive Information
```go
// Good
logger.Info("User authenticated",
    slog.String("user_id", userID),
    slog.String("email_domain", strings.Split(email, "@")[1]),
)

// Avoid
logger.Info("User authenticated",
    slog.String("password", password), // Never log passwords
    slog.String("jwt_token", token),   // Don't log tokens
)
```

## Integration with Middleware

The logger can be integrated with HTTP middleware for automatic request logging:

```go
func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            requestID := generateRequestID()
            
            // Add request ID to context
            ctx := context.WithValue(r.Context(), "request_id", requestID)
            r = r.WithContext(ctx)
            
            // Create request logger
            requestLogger := core.LoggerWithContext(logger, requestID)
            
            // Log request start
            requestLogger.Info("Request started",
                slog.String("method", r.Method),
                slog.String("path", r.URL.Path),
                slog.String("user_agent", r.UserAgent()),
            )
            
            next.ServeHTTP(w, r)
            
            // Log request completion
            requestLogger.Info("Request completed",
                slog.Duration("duration", time.Since(start)),
            )
        })
    }
}
```

## Monitoring and Observability

The structured logging format makes it easy to integrate with monitoring tools:

- **ELK Stack**: JSON format works directly with Elasticsearch
- **Prometheus**: Extract metrics from structured logs
- **Grafana**: Create dashboards from log data
- **AWS CloudWatch**: Parse structured logs for insights

## Performance Considerations

- Use appropriate log levels to avoid excessive logging in production
- The JSON format has slightly more overhead than text format
- Source information (`AddSource: true`) adds some performance cost
- Consider using sampling for high-frequency debug logs

## Migration from Old Logger

To migrate from the old logger configuration:

1. Replace direct `slog.New()` calls with `core.NewDefaultLogger()`
2. Add contextual information using the helper functions
3. Use structured attributes instead of string formatting
4. Set appropriate environment variables for your deployment

## Testing

For testing, you can create a logger that writes to a buffer:

```go
func TestLogger() *slog.Logger {
    var buf bytes.Buffer
    config := &core.LoggerConfig{
        Level:     core.LogLevelDebug,
        Format:    core.LogFormatText,
        AddSource: false,
        Output:    &buf,
    }
    return core.NewLogger(config)
}
```
