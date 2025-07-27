package api

import (
	"net/http"
	"time"

	"github.com/othavioBF/pandoragym-go-api/internal/utils"
)

func (api *API) HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "pandoragym-api",
		"version":   "1.0.0",
	})
}

func (api *API) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "No file provided")
		return
	}
	defer file.Close()

	// TODO: Implement file upload to storage service (Supabase, S3, etc.)
	// For now, just return a mock URL
	fileURL := "https://example.com/uploads/" + header.Filename

	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"file_url": fileURL,
	})
}

func (api *API) GetStatistics(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement statistics gathering
	stats := map[string]any{
		"total_users":     0,
		"total_workouts":  0,
		"total_exercises": 0,
		"active_sessions": 0,
	}

	utils.WriteJSONResponse(w, http.StatusOK, stats)
}
