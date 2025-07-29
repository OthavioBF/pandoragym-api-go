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

// JWTMiddleware validates JWT tokens and adds user info to context
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

		// Get JWT secret from environment
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "your-super-secret-jwt-key-change-this-in-production"
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			api.Logger.Error("Invalid JWT token", "error", err)
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid user ID in token")
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid user ID format")
			return
		}

		// Get user from database
		user, err := api.UserService.GetUserByID(r.Context(), userID)
		if err != nil {
			api.Logger.Error("Failed to get user from token", "error", err, "user_id", userID)
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "User not found")
			return
		}

		if user == nil {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, "User not found")
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), utils.UserIDKey, userID)
		ctx = context.WithValue(ctx, utils.UserRoleKey, user.Role)
		ctx = context.WithValue(ctx, utils.UserKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// PersonalOnlyMiddleware restricts access to personal trainers only
func (api *API) PersonalOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(utils.UserRoleKey).(pgstore.Role)
		if !ok || role != pgstore.RolePersonal {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied: Personal trainer role required")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// StudentOnlyMiddleware restricts access to students only
func (api *API) StudentOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(utils.UserRoleKey).(pgstore.Role)
		if !ok || role != pgstore.RoleStudent {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied: Student role required")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// AdminOnlyMiddleware restricts access to admins only
func (api *API) AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(utils.UserRoleKey).(pgstore.Role)
		if !ok || role != pgstore.RoleAdmin {
			utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied: Admin role required")
			return
		}

		next.ServeHTTP(w, r)
	})
}
