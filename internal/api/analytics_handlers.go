package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/othavioBF/pandoragym-go-api/internal/utils"
)

func (api *API) GetWorkoutFrequency(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	frequency, err := api.AnalyticsService.GetWorkoutFrequency(r.Context(), userID, startDate, endDate)
	if err != nil {
		api.Logger.Error("Failed to get workout frequency", "error", err, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get workout frequency")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, frequency)
}

func (api *API) GetWorkoutHistoryExercises(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	history, err := api.AnalyticsService.GetWorkoutHistoryExercises(r.Context(), userID)
	if err != nil {
		api.Logger.Error("Failed to get workout history exercises", "error", err, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get workout history")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, history)
}

func (api *API) GetWorkoutExercisePerformanceComparison(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	exerciseIDStr := r.URL.Query().Get("exercise_id")
	if exerciseIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "exercise_id parameter is required")
		return
	}

	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid exercise ID format")
		return
	}

	comparison, err := api.AnalyticsService.GetExercisePerformanceComparison(r.Context(), userID, exerciseID)
	if err != nil {
		api.Logger.Error("Failed to get exercise performance comparison", "error", err, "user_id", userID, "exercise_id", exerciseID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get performance comparison")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, comparison)
}

func (api *API) GetWorkoutFrequencyForUser(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userIDStr := chi.URLParam(r, "userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	// TODO: Verify trainer has access to this user
	// hasAccess, err := api.UserService.TrainerHasAccessToStudent(trainerID.String(), userID.String())
	// if err != nil || !hasAccess {
	//     utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
	//     return
	// }

	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	frequency, err := api.AnalyticsService.GetWorkoutFrequency(r.Context(), userID, startDate, endDate)
	if err != nil {
		api.Logger.Error("Failed to get workout frequency for user", "error", err, "trainer_id", trainerID, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get workout frequency")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, frequency)
}

func (api *API) GetWorkoutHistoryForUser(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userIDStr := chi.URLParam(r, "userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	// TODO: Verify trainer has access to this user
	// hasAccess, err := api.UserService.TrainerHasAccessToStudent(trainerID.String(), userID.String())
	// if err != nil || !hasAccess {
	//     utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
	//     return
	// }

	history, err := api.AnalyticsService.GetWorkoutHistoryExercises(r.Context(), userID)
	if err != nil {
		api.Logger.Error("Failed to get workout history for user", "error", err, "trainer_id", trainerID, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get workout history")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, history)
}

func (api *API) GetWorkoutPerformanceForUser(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userIDStr := chi.URLParam(r, "userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	exerciseIDStr := r.URL.Query().Get("exercise_id")
	if exerciseIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "exercise_id parameter is required")
		return
	}

	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid exercise ID format")
		return
	}

	// TODO: Verify trainer has access to this user
	// hasAccess, err := api.UserService.TrainerHasAccessToStudent(trainerID.String(), userID.String())
	// if err != nil || !hasAccess {
	//     utils.WriteErrorResponse(w, http.StatusForbidden, "Access denied")
	//     return
	// }

	performance, err := api.AnalyticsService.GetExercisePerformanceComparison(r.Context(), userID, exerciseID)
	if err != nil {
		api.Logger.Error("Failed to get workout performance for user", "error", err, "trainer_id", trainerID, "user_id", userID, "exercise_id", exerciseID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get performance data")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, performance)
}
