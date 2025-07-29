package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/othavioBF/pandoragym-go-api/internal/utils"
)

// Training program handlers (now using workout service)

func (api *API) GetAllTrainingPrograms(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	programs, err := api.WorkoutService.GetAllPrograms(r.Context(), userID.String())
	if err != nil {
		api.Logger.Error("Failed to get all training programs", "error", err, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get training programs")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, programs)
}

func (api *API) GetFreeTrainingPrograms(w http.ResponseWriter, r *http.Request) {
	programs, err := api.WorkoutService.GetFreePrograms(r.Context())
	if err != nil {
		api.Logger.Error("Failed to get free training programs", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get free training programs")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"programs": programs,
	})
}

func (api *API) GetFreeTrainingProgramByID(w http.ResponseWriter, r *http.Request) {
	programID := chi.URLParam(r, "id")

	program, err := api.WorkoutService.GetFreeProgramByID(r.Context(), programID)
	if err != nil {
		api.Logger.Error("Failed to get free training program", "error", err, "program_id", programID)
		utils.WriteErrorResponse(w, http.StatusNotFound, "Training program not found")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, program)
}

// Subscription management handlers

func (api *API) SubscribeToPlan(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	req, err := utils.DecodeValidJSON[struct {
		PlanID string `json:"plan_id"`
	}](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// For now, just return success
	// TODO: Implement actual subscription logic
	api.Logger.Info("Plan subscription requested", "user_id", userID, "plan_id", req.PlanID)

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Subscribed to plan successfully",
	})
}

func (api *API) CancelPlan(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// For now, just return success
	// TODO: Implement actual plan cancellation logic
	api.Logger.Info("Plan cancellation requested", "user_id", userID)

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Plan cancelled successfully",
	})
}
