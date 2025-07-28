package core

import (
	"context"
	"log/slog"
	"net/http"
)

// This file contains examples of how to use the improved logger
// These are example functions to demonstrate usage patterns

// Example: Using logger with request context in HTTP handlers
func ExampleHTTPHandler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create logger with request context
		requestLogger := LoggerWithContext(logger, "req-123")
		
		// Log request start
		// Output:
		// 2025-07-28 19:00:00 INFO:
		//     [handlers.go:15] Processing request method=GET path=/api/users user_agent=Mozilla/5.0
		requestLogger.Info("Processing request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("user_agent", r.UserAgent()),
		)
		
		// Your handler logic here...
		
		// Log successful response
		// Output:
		// 2025-07-28 19:00:00 INFO:
		//     [handlers.go:25] Request completed successfully status=200
		requestLogger.Info("Request completed successfully",
			slog.Int("status", http.StatusOK),
		)
	}
}

// Example: Using logger in service layer
func ExampleServiceMethod(logger *slog.Logger, userID string) error {
	// Create service-specific logger
	serviceLogger := LoggerWithService(logger, "user-service")
	
	// Add user context
	userLogger := LoggerWithUser(serviceLogger, userID)
	
	// Output:
	// 2025-07-28 19:00:00 DEBUG:
	//     [user_service.go:42] Starting user operation service=user-service user_id=123
	userLogger.Debug("Starting user operation")
	
	// Simulate some work
	if userID == "" {
		// Output:
		// 2025-07-28 19:00:00 ERROR:
		//     [user_service.go:48] Invalid user ID provided validation_error=user_id_empty service=user-service user_id=
		userLogger.Error("Invalid user ID provided",
			slog.String("validation_error", "user_id_empty"),
		)
		return nil
	}
	
	// Output:
	// 2025-07-28 19:00:00 INFO:
	//     [user_service.go:55] User operation completed successfully operation=update_profile service=user-service user_id=123
	userLogger.Info("User operation completed successfully",
		slog.String("operation", "update_profile"),
	)
	
	return nil
}

// Example: Using structured logging for database operations
func ExampleDatabaseOperation(logger *slog.Logger, ctx context.Context) {
	dbLogger := logger.With(
		slog.String("component", "database"),
		slog.String("operation", "user_query"),
	)
	
	// Output:
	// 2025-07-28 19:00:00 DEBUG:
	//     [db.go:25] Executing database query query=SELECT * FROM users WHERE id = $1 component=database operation=user_query
	dbLogger.Debug("Executing database query",
		slog.String("query", "SELECT * FROM users WHERE id = $1"),
	)
	
	// Simulate database operation
	// ...
	
	// Output:
	// 2025-07-28 19:00:00 INFO:
	//     [db.go:32] Database query completed rows_affected=1 duration=0s component=database operation=user_query
	dbLogger.Info("Database query completed",
		slog.Int("rows_affected", 1),
		slog.Duration("duration", 0), // You would measure actual duration
	)
}

// Example: Error logging with stack trace context
func ExampleErrorLogging(logger *slog.Logger, err error) {
	// Output:
	// 2025-07-28 19:00:00 ERROR:
	//     [user_service.go:42] Operation failed error=connection timeout error_type=validation_error context=[function=CreateUser file=user_service.go:42]
	logger.Error("Operation failed",
		slog.String("error", err.Error()),
		slog.String("error_type", "validation_error"),
		slog.Group("context",
			slog.String("function", "CreateUser"),
			slog.String("file", "user_service.go:42"),
		),
	)
}

// Example: Performance monitoring with structured logs
func ExamplePerformanceLogging(logger *slog.Logger) {
	perfLogger := logger.With(
		slog.String("component", "performance"),
	)
	
	// Output:
	// 2025-07-28 19:00:00 INFO:
	//     [monitor.go:15] Performance metrics metrics=[active_connections=25 requests_per_second=150 avg_response_time_ms=45.2 memory_usage_mb=128] component=performance
	perfLogger.Info("Performance metrics",
		slog.Group("metrics",
			slog.Int("active_connections", 25),
			slog.Int("requests_per_second", 150),
			slog.Float64("avg_response_time_ms", 45.2),
			slog.Int64("memory_usage_mb", 128),
		),
	)
}

// Example: Business logic logging
func ExampleBusinessLogic(logger *slog.Logger, userID, workoutID string) {
	businessLogger := logger.With(
		slog.String("domain", "workout"),
		slog.String("action", "assign_workout"),
	)
	
	// Output:
	// 2025-07-28 19:00:00 INFO:
	//     [workout_service.go:25] Assigning workout to user user_id=123 workout_id=456 assigned_by=personal_trainer domain=workout action=assign_workout
	businessLogger.Info("Assigning workout to user",
		slog.String("user_id", userID),
		slog.String("workout_id", workoutID),
		slog.String("assigned_by", "personal_trainer"),
	)
	
	// Business logic here...
	
	// Output:
	// 2025-07-28 19:00:00 INFO:
	//     [workout_service.go:35] Workout assigned successfully user_id=123 workout_id=456 notification_sent=true domain=workout action=assign_workout
	businessLogger.Info("Workout assigned successfully",
		slog.String("user_id", userID),
		slog.String("workout_id", workoutID),
		slog.Bool("notification_sent", true),
	)
}
