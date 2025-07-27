package core

import (
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/othavioBF/pandoragym-go-api/internal/api"
	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
	"github.com/othavioBF/pandoragym-go-api/internal/services"
)

func InjectDependencies(queries *pgstore.Queries) api.API {
	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				AddSource: true,
			},
		),
	)

	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key-change-this-in-production"
	}

	userService := services.NewUserService(queries)
	workoutService := services.NewWorkoutService(queries)
	schedulingService := services.NewSchedulingService(queries)
	authService := services.NewAuthService(queries, jwtSecret)

	return api.API{
		Router:            chi.NewMux(),
		Logger:            logger,
		UserService:       userService,
		WorkoutService:    workoutService,
		SchedulingService: schedulingService,
		AuthService:       authService,
	}
}
