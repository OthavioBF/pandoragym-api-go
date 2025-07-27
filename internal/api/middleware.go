package api

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
	"github.com/othavioBF/pandoragym-go-api/internal/utils"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
)

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// JWTMiddleware is attached to the API struct to access services if needed
func (api *API) JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Bearer token required")
			return
		}

		// Get JWT secret from environment variable
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "your-secret-key" // Fallback for development
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		// Parse user ID from token claims
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid user ID in token")
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid user ID format in token")
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// PersonalOnlyMiddleware restricts access to personal trainers only
func (api *API) PersonalOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := api.GetUserIDFromContext(r.Context())
		if userID == uuid.Nil {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Check if user is a personal trainer by querying the database
		user, err := api.UserService.GetUserByID(r.Context(), userID)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		if user.Role != pgstore.RolePersonal {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied: Personal trainer role required")
			return
		}

		// Add user role to context for future use
		ctx := context.WithValue(r.Context(), UserRoleKey, user.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *API) StudentOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := api.GetUserIDFromContext(r.Context())
		if userID == uuid.Nil {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Check if user is a student by querying the database
		user, err := api.UserService.GetUserByID(r.Context(), userID)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		if user.Role != pgstore.RoleStudent {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied: Student role required")
			return
		}

		ctx := context.WithValue(r.Context(), UserRoleKey, user.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminOnlyMiddleware restricts access to admins only
func (api *API) AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := api.GetUserIDFromContext(r.Context())
		if userID == uuid.Nil {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		user, err := api.UserService.GetUserByID(r.Context(), userID)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		if user.Role != pgstore.RolePersonal { // Temporary: treat personal trainers as admins
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied: Admin role required")
			return
		}

		ctx := context.WithValue(r.Context(), UserRoleKey, user.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *API) GetUserIDFromContext(ctx context.Context) uuid.UUID {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return userID
}

func (api *API) GetUserRoleFromContext(ctx context.Context) pgstore.Role {
	role, ok := ctx.Value(UserRoleKey).(pgstore.Role)
	if !ok {
		return ""
	}
	return role
}
