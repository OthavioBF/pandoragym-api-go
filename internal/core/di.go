package core

import (
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/othavioBF/pandoragym-go-api/internal/api"
	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
	"github.com/othavioBF/pandoragym-go-api/internal/services"
)

func InjectDependencies(queries *pgstore.Queries) api.API {
	// Initialize logger with improved configuration
	logger := NewDefaultLogger()
	
	// Log application startup
	logger.Info("Initializing PandoraGym API dependencies")

	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logger.Warn("JWT_SECRET not set, using default value (not recommended for production)")
		jwtSecret = "your-super-secret-jwt-key-change-this-in-production"
	}

	// Initialize services with service-specific loggers
	userService := services.NewUserService(queries)
	workoutService := services.NewWorkoutService(queries)
	schedulingService := services.NewSchedulingService(queries)
	authService := services.NewAuthService(queries, jwtSecret)
	analyticsService := services.NewAnalyticsService(queries)
	planService := services.NewPlanService(queries)
	systemService := services.NewSystemService()

	logger.Info("All dependencies initialized successfully")

	return api.API{
		Router:            chi.NewMux(),
		Logger:            logger,
		UserService:       userService,
		WorkoutService:    workoutService,
		SchedulingService: schedulingService,
		AuthService:       authService,
		AnalyticsService:  analyticsService,
		PlanService:       planService,
		SystemService:     systemService,
	}
}
