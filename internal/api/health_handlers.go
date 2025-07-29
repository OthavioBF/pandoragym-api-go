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
