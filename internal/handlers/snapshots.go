package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/errors"
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
)

func (h *Handler) SnapshotExport(w http.ResponseWriter, r *http.Request) error {
	var req config.SnapshotExportRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.BadRequest("Invalid request body", "")
	}
	defer r.Body.Close()

	if req.ProjectName == "" {
		return errors.Validation(map[string]string{
			"project_name": "project_name is required",
		})
	}

	res, err := h.Services.Snapshot.ExportSnapshot(r.Context(), req)
	if err != nil {
		return err
	}

	helpers.WriteResponse(w, http.StatusOK, res)
	return nil
}

func (h *Handler) SnapshotImport(w http.ResponseWriter, r *http.Request) error {
	var req config.SnapshotImportRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errors.BadRequest("Invalid request body", "")
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
		return errors.Validation(validationErrors)
	}

	res, err := h.Services.Snapshot.ImportSnapshot(r.Context(), req)
	if err != nil {
		return err
	}

	helpers.WriteResponse(w, http.StatusCreated, res)
	return nil
}
