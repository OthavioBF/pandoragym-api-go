package core

import (
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/othavioBF/pandoragym-go-api/internal/api"
	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
	"github.com/othavioBF/pandoragym-go-api/internal/services"
)

func InjectDependencies(queries *pgstore.Queries, pool *pgxpool.Pool) api.API {
	logger := NewDefaultLogger()

	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(pool)
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Name = "pandoragym_session"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.SameSite = 3
	sessionManager.Cookie.Secure = os.Getenv("ENV") == "production"

	userService := services.NewUserService(queries)
	workoutService := services.NewWorkoutService(queries, pool)
	schedulingService := services.NewSchedulingService(queries)
	authService := services.NewAuthService(queries, sessionManager)
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
		SessionManager:    sessionManager,
	}
}
