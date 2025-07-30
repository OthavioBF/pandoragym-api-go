package api

import (
	"log/slog"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/othavioBF/pandoragym-go-api/internal/services"
)

type API struct {
	Router            *chi.Mux
	Logger            *slog.Logger
	UserService       *services.UserService
	WorkoutService    *services.WorkoutService
	SchedulingService *services.SchedulingService
	AuthService       *services.AuthService
	AnalyticsService  *services.AnalyticsService
	PlanService       *services.PlanService
	SystemService     services.SystemService
	SessionManager    *scs.SessionManager
}
