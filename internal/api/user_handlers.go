package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/othavioBF/pandoragym-go-api/internal/middleware/auth"
	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
	"github.com/othavioBF/pandoragym-go-api/internal/utils"
)

func (api *API) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := api.UserService.GetUserByID(r.Context(), userID)
	if err != nil {
		api.Logger.Error("Failed to get user profile", "error", err, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"user": user,
	})
}

func (api *API) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	req, err := utils.DecodeValidJSON[pgstore.UpdateProfileRequest](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = api.UserService.UpdateUserProfile(r.Context(), userID, &req)
	if err != nil {
		api.Logger.Error("Failed to update user profile", "error", err, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	// Get updated user
	user, err := api.UserService.GetUserByID(r.Context(), userID)
	if err != nil {
		api.Logger.Error("Failed to get updated user profile", "error", err, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get updated profile")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"user": user,
	})
}

func (api *API) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "No file provided")
		return
	}
	defer file.Close()

	// Validate file type
	if !utils.IsValidImageType(header.Header.Get("Content-Type")) {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid file type. Only images are allowed")
		return
	}

	// For now, just return a placeholder URL
	// TODO: Implement actual file upload to Supabase/S3
	avatarURL := "https://example.com/avatars/" + userID.String() + ".jpg"

	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"avatar_url": avatarURL,
	})
}

func (api *API) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// This is an admin-only endpoint
	// TODO: Add admin check
	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"users": []pgstore.User{}, // Placeholder
	})
}

// Personal trainer specific handlers
func (api *API) GetPersonalStudents(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// TODO: Implement getting personal trainer's students
	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"students": []pgstore.User{},
	})
}

func (api *API) GetStudentEvolution(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// TODO: Implement student evolution tracking
	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"evolution": map[string]interface{}{},
	})
}

func (api *API) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// TODO: Implement messaging system
	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"message": "Message sent successfully",
	})
}

func (api *API) GetPersonalSchedule(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// TODO: Implement personal schedule management
	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"schedule": []pgstore.PersonalSchedule{},
	})
}

func (api *API) CreatePersonalSchedule(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// TODO: Implement personal schedule creation
	utils.WriteJSONResponse(w, http.StatusCreated, map[string]any{
		"message": "Schedule created successfully",
	})
}

// Student specific handlers
func (api *API) GetWorkoutHistory(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// TODO: Implement workout history
	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"history": []pgstore.WorkoutsHistory{},
	})
}

func (api *API) ExecuteWorkout(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// TODO: Implement workout execution tracking
	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"message": "Workout execution recorded",
	})
}

func (api *API) RateWorkout(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// TODO: Implement workout rating system
	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"message": "Workout rated successfully",
	})
}
