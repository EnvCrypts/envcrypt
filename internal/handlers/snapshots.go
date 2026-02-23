package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
)

func (h *Handler) SnapshotExport(w http.ResponseWriter, r *http.Request) {
	var req config.SnapshotExportRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	defer r.Body.Close()

	if req.ProjectName == "" {
		helpers.WriteError(w, http.StatusBadRequest, helpers.ErrValidation(map[string]string{
			"project_name": "project_name is required",
		}))
		return
	}

	res, err := h.Services.Snapshot.ExportSnapshot(r.Context(), req)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	helpers.WriteResponse(w, http.StatusOK, res)
}

func (h *Handler) SnapshotImport(w http.ResponseWriter, r *http.Request) {
	var req config.SnapshotImportRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	defer r.Body.Close()

	validationErrors := make(map[string]string)
	if req.NewProjectName == "" {
		validationErrors["new_project_name"] = "new_project_name is required"
	}
	if req.Checksum == "" {
		validationErrors["checksum"] = "checksum is required"
	}
	if len(validationErrors) > 0 {
		helpers.WriteError(w, http.StatusBadRequest, helpers.ErrValidation(validationErrors))
		return
	}

	res, err := h.Services.Snapshot.ImportSnapshot(r.Context(), req)
	if err != nil {
		helpers.WriteError(w, 0, err)
		return
	}

	helpers.WriteResponse(w, http.StatusCreated, res)
}
