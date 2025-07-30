package core

import (
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/othavioBF/pandoragym-go-api/internal/api"
	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
	"github.com/othavioBF/pandoragym-go-api/internal/services"
)

func InjectDependencies(queries *pgstore.Queries, pool *pgxpool.Pool) api.API {
	logger := NewDefaultLogger()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logger.Warn("JWT_SECRET not set, using default value (not recommended for production)")
		jwtSecret = "your-super-secret-jwt-key-change-this-in-production"
	}

	userService := services.NewUserService(queries)
	workoutService := services.NewWorkoutService(queries, pool)
	schedulingService := services.NewSchedulingService(queries)
	authService := services.NewAuthService(queries, jwtSecret)
	analyticsService := services.NewAnalyticsService(queries)
	planService := services.NewPlanService(queries)
	systemService := services.NewSystemService()

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
