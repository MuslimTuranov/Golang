package handlers

import (
	"net/http"
)

// Health godoc
// @Summary Healthcheck
// @Description Returns API health status
// @Tags health
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} healthResponse
// @Router /healthz [get]
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}
