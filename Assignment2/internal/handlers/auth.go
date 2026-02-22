package handlers

import (
	"Assignment2/pkg/modules"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type loginReq struct {
	UserID   int    `json:"user_id"`
	Password string `json:"password"`
}

// Login godoc
// @Summary Login
// @Description Authenticates user and returns JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body loginReq true "Login payload"
// @Success 200 {object} tokenResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /auth/login [post]
func (h *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		token, err := h.Auth.Login(ctx, req.UserID, req.Password)
		if err != nil {
			if errors.Is(err, modules.ErrUnauthorized) {
				writeError(w, http.StatusUnauthorized, "invalid credentials")
				return
			}
			writeError(w, http.StatusInternalServerError, "login failed")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"token": token})
	}
}
