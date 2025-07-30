package api

import (
	"context"
	"net/http"

	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
	"github.com/othavioBF/pandoragym-go-api/internal/utils"
)

func (api *API) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !api.AuthService.IsAuthenticated(r.Context()) {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		userID, err := api.AuthService.GetUserIDFromSession(r.Context())
		if err != nil {
			api.Logger.Error("Failed to get user ID from session", "error", err)
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid session")
			return
		}

		ctx := context.WithValue(r.Context(), utils.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *API) RequireRole(roles ...pgstore.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, err := api.AuthService.GetUserRoleFromSession(r.Context())
			if err != nil {
				api.Logger.Error("Failed to get user role from session", "error", err)
				utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid session")
				return
			}

			role := pgstore.Role(userRole)
			hasPermission := false
			for _, allowedRole := range roles {
				if role == allowedRole {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				utils.WriteErrorResponse(w, http.StatusForbidden, "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (api *API) RequireStudent(next http.Handler) http.Handler {
	return api.RequireRole(pgstore.RoleStudent)(next)
}

func (api *API) RequirePersonal(next http.Handler) http.Handler {
	return api.RequireRole(pgstore.RolePersonal)(next)
}

func (api *API) RequireAdmin(next http.Handler) http.Handler {
	return api.RequireRole(pgstore.RoleAdmin)(next)
}

func (api *API) RequireStudentOrPersonal(next http.Handler) http.Handler {
	return api.RequireRole(pgstore.RoleStudent, pgstore.RolePersonal)(next)
}

func (api *API) RequirePersonalOrAdmin(next http.Handler) http.Handler {
	return api.RequireRole(pgstore.RolePersonal, pgstore.RoleAdmin)(next)
}

func (api *API) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (api *API) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.Logger.Info("HTTP Request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)

		next.ServeHTTP(w, r)
	})
}
