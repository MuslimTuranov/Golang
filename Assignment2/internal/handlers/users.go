package handlers

import (
	"Assignment2/pkg/modules"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// GetUsers godoc
// @Summary List users
// @Description Returns non-deleted users with pagination
// @Tags users
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param limit query int false "Page size" minimum(1) maximum(200)
// @Param offset query int false "Rows offset" minimum(0)
// @Success 200 {array} modules.User
// @Failure 500 {object} errorResponse
// @Router /users/ [get]
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	users, err := h.Users.GetUsers(ctx, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch users")
		return
	}
	writeJSON(w, http.StatusOK, users)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Tags users
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} modules.User
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /users/{id} [get]
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	u, err := h.Users.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, modules.ErrNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}
	writeJSON(w, http.StatusOK, u)
}

type createUserReq struct {
	Name     string  `json:"name"`
	Email    *string `json:"email"`
	Age      *int    `json:"age"`
	Password string  `json:"password"`
}

// CreateUser godoc
// @Summary Create user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param request body createUserReq true "Create user payload"
// @Success 201 {object} idResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /users/ [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id, err := h.Users.CreateUser(ctx, req.Name, req.Email, req.Age, req.Password)
	if err != nil {
		if errors.Is(err, modules.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, "invalid input")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create user")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"id": id})
}

type updateUserReq struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Age      *int    `json:"age"`
	Password *string `json:"password"`
}

// UpdateUser godoc
// @Summary Update user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body updateUserReq true "Update user payload"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /users/{id} [put]
// @Router /users/{id} [patch]
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req updateUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.Users.UpdateUser(ctx, id, req.Name, req.Email, req.Age, req.Password); err != nil {
		if errors.Is(err, modules.ErrNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update user")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "updated"})
}

// DeleteUser godoc
// @Summary Soft delete user
// @Tags users
// @Produce json
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} deleteResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /users/{id} [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	ra, err := h.Users.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, modules.ErrNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete user")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"rows_affected": ra})
}
