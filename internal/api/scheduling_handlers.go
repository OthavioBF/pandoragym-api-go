package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/othavioBF/pandoragym-go-api/internal/middleware/auth"
	"github.com/othavioBF/pandoragym-go-api/internal/infra/pgstore"
	"github.com/othavioBF/pandoragym-go-api/internal/utils"
)

func (api *API) GetSchedulings(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
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
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	schedulingIDStr := chi.URLParam(r, "id")
	if schedulingIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Scheduling ID is required")
		return
	}

	schedulingID, err := uuid.Parse(schedulingIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid scheduling ID format")
		return
	}

	scheduling, err := api.SchedulingService.GetSchedulingByID(r.Context(), schedulingID, userID)
	if err != nil {
		api.Logger.Error("Failed to get scheduling", "error", err, "scheduling_id", schedulingID, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusNotFound, "Scheduling not found")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"scheduling": scheduling,
	})
}

func (api *API) CreateScheduling(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	req, err := utils.DecodeValidJSON[pgstore.CreateSchedulingParams](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	scheduling, err := api.SchedulingService.CreateScheduling(r.Context(), userID, req)
	if err != nil {
		api.Logger.Error("Failed to create scheduling", "error", err, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create scheduling")
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, map[string]any{
		"scheduling": scheduling,
	})
}

func (api *API) UpdateScheduling(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	schedulingIDStr := chi.URLParam(r, "id")
	if schedulingIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Scheduling ID is required")
		return
	}

	schedulingID, err := uuid.Parse(schedulingIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid scheduling ID format")
		return
	}

	req, err := utils.DecodeValidJSON[pgstore.UpdateSchedulingRequest](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = api.SchedulingService.UpdateScheduling(r.Context(), schedulingID, userID, req)
	if err != nil {
		api.Logger.Error("Failed to update scheduling", "error", err, "scheduling_id", schedulingID, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update scheduling")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"message": "Scheduling updated successfully",
	})
}

func (api *API) CancelScheduling(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	schedulingIDStr := chi.URLParam(r, "id")
	if schedulingIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Scheduling ID is required")
		return
	}

	schedulingID, err := uuid.Parse(schedulingIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid scheduling ID format")
		return
	}

	req, err := utils.DecodeValidJSON[pgstore.CancelSchedulingRequest](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = api.SchedulingService.CancelScheduling(r.Context(), schedulingID, userID, req.Reason)
	if err != nil {
		api.Logger.Error("Failed to cancel scheduling", "error", err, "scheduling_id", schedulingID, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to cancel scheduling")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"message": "Scheduling canceled successfully",
	})
}

func (api *API) StartScheduling(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	schedulingIDStr := chi.URLParam(r, "id")
	if schedulingIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Scheduling ID is required")
		return
	}

	schedulingID, err := uuid.Parse(schedulingIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid scheduling ID format")
		return
	}

	err = api.SchedulingService.StartScheduling(r.Context(), schedulingID, userID)
	if err != nil {
		api.Logger.Error("Failed to start scheduling", "error", err, "scheduling_id", schedulingID, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to start scheduling")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"message": "Scheduling started successfully",
	})
}

func (api *API) CompleteScheduling(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID == uuid.Nil {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	schedulingIDStr := chi.URLParam(r, "id")
	if schedulingIDStr == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Scheduling ID is required")
		return
	}

	schedulingID, err := uuid.Parse(schedulingIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid scheduling ID format")
		return
	}

	err = api.SchedulingService.CompleteScheduling(r.Context(), schedulingID, userID)
	if err != nil {
		api.Logger.Error("Failed to complete scheduling", "error", err, "scheduling_id", schedulingID, "user_id", userID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to complete scheduling")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"message": "Scheduling completed successfully",
	})
}
