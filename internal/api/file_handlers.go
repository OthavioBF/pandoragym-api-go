package api

import (
	"net/http"

	"github.com/othavioBF/pandoragym-go-api/internal/utils"
)

func (api *API) UploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
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

	url, err := api.FileService.UploadFile(r.Context(), file, header)
	if err != nil {
		api.Logger.Error("Failed to get personal trainer", "error", err, "trainer_id", "")
		utils.WriteErrorResponse(w, http.StatusNotFound, "Personal trainer not found")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]any{
		"message": "File uploaded successfully",
		"url":     url,
	})
}

func (api *API) DeleteFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "No file provided")
		return
	}
	defer file.Close()

	err = api.FileService.DeleteFile(r.Context(), "")
	if err != nil {
		api.Logger.Error("Failed to get personal trainer", "error", err, "trainer_id", "")
		utils.WriteErrorResponse(w, http.StatusNotFound, "Personal trainer not found")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Arquivo deletado com sucesso",
	})
}
