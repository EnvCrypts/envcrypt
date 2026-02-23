package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
)

func (h *Handler) SnapshotExport(w http.ResponseWriter, r *http.Request) {
	var req config.SnapshotExportRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if req.ProjectName == "" {
		helpers.WriteError(w, http.StatusBadRequest, "project_name is required")
		return
	}

	res, err := h.Services.Snapshot.ExportSnapshot(r.Context(), req)
	if err != nil {
		if err.Error() == "project not found or permission denied" || err.Error() == "permission denied" {
			helpers.WriteError(w, http.StatusForbidden, err.Error())
			return
		}
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.WriteResponse(w, http.StatusOK, res)
}

func (h *Handler) SnapshotImport(w http.ResponseWriter, r *http.Request) {
	var req config.SnapshotImportRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helpers.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if req.NewProjectName == "" {
		helpers.WriteError(w, http.StatusBadRequest, "new_project_name is required")
		return
	}
	if req.Checksum == "" {
		helpers.WriteError(w, http.StatusBadRequest, "checksum is required")
		return
	}

	res, err := h.Services.Snapshot.ImportSnapshot(r.Context(), req)
	if err != nil {
		if err.Error() == "checksum mismatch" {
			helpers.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err.Error() == "project with this name already exists" {
			helpers.WriteError(w, http.StatusConflict, err.Error())
			return
		}
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.WriteResponse(w, http.StatusCreated, res)
}
