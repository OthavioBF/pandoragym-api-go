package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/othavioBF/pandoragym-go-api/internal/utils"
)

func (api *API) GetTrainerPlans(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	plans, err := api.PlanService.GetTrainerPlans(r.Context(), trainerID)
	if err != nil {
		api.Logger.Error("Failed to get trainer plans", "error", err, "trainer_id", trainerID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get plans")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"plans": plans,
	})
}

func (api *API) CreatePlan(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	req, err := utils.DecodeValidJSON[struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Price       float64  `json:"price"`
		Duration    int      `json:"duration"` // in days
		Features    []string `json:"features"`
	}](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	plan, err := api.PlanService.CreatePlan(r.Context(), trainerID, req.Name, req.Description, req.Price, req.Duration, req.Features)
	if err != nil {
		api.Logger.Error("Failed to create plan", "error", err, "trainer_id", trainerID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create plan")
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, plan)
}

func (api *API) UpdatePlan(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	planID := chi.URLParam(r, "id")

	req, err := utils.DecodeValidJSON[struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Price       float64  `json:"price"`
		Duration    int      `json:"duration"`
		Features    []string `json:"features"`
	}](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	plan, err := api.PlanService.UpdatePlan(r.Context(), trainerID, planID, req.Name, req.Description, req.Price, req.Duration, req.Features)
	if err != nil {
		api.Logger.Error("Failed to update plan", "error", err, "trainer_id", trainerID, "plan_id", planID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update plan")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, plan)
}

func (api *API) DeletePlan(w http.ResponseWriter, r *http.Request) {
	trainerID, ok := r.Context().Value(utils.UserIDKey).(uuid.UUID)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	planID := chi.URLParam(r, "id")

	err := api.PlanService.DeletePlan(r.Context(), trainerID, planID)
	if err != nil {
		api.Logger.Error("Failed to delete plan", "error", err, "trainer_id", trainerID, "plan_id", planID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete plan")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Plan deleted successfully",
	})
}

func (api *API) SubscribeToTrainerPlan(w http.ResponseWriter, r *http.Request) {
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

func (api *API) CancelTrainerPlan(w http.ResponseWriter, r *http.Request) {
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
