package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/othavioBF/pandoragym-go-api/internal/utils"
)

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

func (api *API) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
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
