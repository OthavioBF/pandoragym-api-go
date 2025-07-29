package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/othavioBF/pandoragym-go-api/internal/utils"
)

// Platform management (admin only)

func (api *API) GetStatistics(w http.ResponseWriter, r *http.Request) {
	stats, err := api.AnalyticsService.GetPlatformStatistics(r.Context())
	if err != nil {
		api.Logger.Error("Failed to get platform statistics", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get statistics")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, stats)
}

func (api *API) GetReports(w http.ResponseWriter, r *http.Request) {
	reportType := r.URL.Query().Get("type") // users, workouts, revenue, etc.
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	reports, err := api.AnalyticsService.GetReports(r.Context(), reportType, startDate, endDate)
	if err != nil {
		api.Logger.Error("Failed to get reports", "error", err, "type", reportType)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get reports")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, reports)
}

func (api *API) GetSystemHealth(w http.ResponseWriter, r *http.Request) {
	// For now, return basic health status
	// TODO: Implement actual system health checks
	health := map[string]interface{}{
		"status":    "healthy",
		"database":  "connected",
		"services":  "operational",
		"timestamp": "2024-01-01T00:00:00Z",
	}

	utils.WriteJSONResponse(w, http.StatusOK, health)
}

// Template management (admin only) - now using workout service

// Exercise templates
func (api *API) GetExerciseTemplatesAdmin(w http.ResponseWriter, r *http.Request) {
	templates, err := api.WorkoutService.GetAllExerciseTemplates(r.Context())
	if err != nil {
		api.Logger.Error("Failed to get exercise templates", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get exercise templates")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"templates": templates,
	})
}

func (api *API) CreateExerciseTemplate(w http.ResponseWriter, r *http.Request) {
	req, err := utils.DecodeValidJSON[struct {
		Name         string   `json:"name"`
		Description  string   `json:"description"`
		VideoURL     string   `json:"video_url"`
		Instructions string   `json:"instructions"`
		Category     string   `json:"category"`
		MuscleGroups []string `json:"muscle_groups"`
		Equipment    []string `json:"equipment"`
		Difficulty   string   `json:"difficulty"`
	}](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	template, err := api.WorkoutService.CreateExerciseTemplate(
		r.Context(),
		req.Name,
		req.Description,
		req.VideoURL,
		req.Instructions,
		req.Category,
		req.MuscleGroups,
		req.Equipment,
		req.Difficulty,
	)
	if err != nil {
		api.Logger.Error("Failed to create exercise template", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create exercise template")
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, template)
}

func (api *API) DeleteExerciseTemplate(w http.ResponseWriter, r *http.Request) {
	templateID := chi.URLParam(r, "id")

	err := api.WorkoutService.DeleteExerciseTemplate(r.Context(), templateID)
	if err != nil {
		api.Logger.Error("Failed to delete exercise template", "error", err, "template_id", templateID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete exercise template")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Exercise template deleted successfully",
	})
}

// Workout templates
func (api *API) GetWorkoutTemplatesAdmin(w http.ResponseWriter, r *http.Request) {
	templates, err := api.WorkoutService.GetAllWorkoutTemplates(r.Context())
	if err != nil {
		api.Logger.Error("Failed to get workout templates", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get workout templates")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"templates": templates,
	})
}

func (api *API) CreateWorkoutTemplate(w http.ResponseWriter, r *http.Request) {
	req, err := utils.DecodeValidJSON[struct {
		Name        string                   `json:"name"`
		Description string                   `json:"description"`
		Thumbnail   string                   `json:"thumbnail"`
		Category    string                   `json:"category"`
		Difficulty  string                   `json:"difficulty"`
		Duration    int                      `json:"duration"` // in minutes
		WeekDays    []string                 `json:"week_days"`
		Exercises   []map[string]interface{} `json:"exercises"`
		Tags        []string                 `json:"tags"`
	}](r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	template, err := api.WorkoutService.CreateWorkoutTemplate(
		r.Context(),
		req.Name,
		req.Description,
		req.Thumbnail,
		req.Category,
		req.Difficulty,
		req.Duration,
		req.WeekDays,
		req.Exercises,
		req.Tags,
	)
	if err != nil {
		api.Logger.Error("Failed to create workout template", "error", err)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create workout template")
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, template)
}

func (api *API) DeleteWorkoutTemplate(w http.ResponseWriter, r *http.Request) {
	templateID := chi.URLParam(r, "id")

	err := api.WorkoutService.DeleteWorkoutTemplate(r.Context(), templateID)
	if err != nil {
		api.Logger.Error("Failed to delete workout template", "error", err, "template_id", templateID)
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to delete workout template")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Workout template deleted successfully",
	})
}
