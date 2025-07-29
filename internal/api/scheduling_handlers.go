package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
	"github.com/othavioBF/pandoragym-go-api/internal/utils"
)

func (api *API) GetSchedulings(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	schedulings, err := api.SchedulingService.GetSchedulings(r.Context(), userID)
	if err != nil {
		api.Logger.Error("Failed to get schedulings", "error", err, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get schedulings")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"schedulings": schedulings,
	})
}

func (api *API) GetScheduling(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	schedulingIDStr := chi.URLParam(r, "id")
	schedulingID, err := uuid.Parse(schedulingIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid scheduling ID")
		return
	}

	scheduling, err := api.SchedulingService.GetSchedulingByID(r.Context(), schedulingID, userID)
	if err != nil {
		api.Logger.Error("Failed to get scheduling", "error", err, "scheduling_id", schedulingID, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get scheduling")
		return
	}

	if scheduling == nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Scheduling not found")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, scheduling)
}

func (api *API) CreateScheduling(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	req, err := utils.DecodeValidJSON[pgstore.CreateSchedulingRequest](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// For now, return mock response
	// TODO: Implement actual scheduling creation
	api.Logger.Info("Scheduling creation requested", "user_id", userID, "personal_id", req.PersonalID)

	utils.WriteJSONResponse(w, http.StatusCreated, map[string]string{
		"message": "Scheduling created successfully",
		"id":      uuid.New().String(),
	})
}

func (api *API) UpdateScheduling(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	schedulingIDStr := chi.URLParam(r, "id")
	schedulingID, err := uuid.Parse(schedulingIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid scheduling ID")
		return
	}

	req, err := utils.DecodeValidJSON[pgstore.UpdateSchedulingRequest](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// For now, return mock response
	// TODO: Implement actual scheduling update
	api.Logger.Info("Scheduling update requested", "scheduling_id", schedulingID, "user_id", userID, "status", req.Status)

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Scheduling updated successfully",
	})
}

func (api *API) CancelScheduling(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	schedulingIDStr := chi.URLParam(r, "id")
	schedulingID, err := uuid.Parse(schedulingIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid scheduling ID")
		return
	}

	// For now, return mock response
	// TODO: Implement actual scheduling cancellation
	api.Logger.Info("Scheduling cancellation requested", "scheduling_id", schedulingID, "user_id", userID)

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Scheduling cancelled successfully",
	})
}
