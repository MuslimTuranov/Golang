package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"users-service/internal/repository"
)

type UserHandler struct {
	repo *repository.Repository
}

func NewUserHandler(repo *repository.Repository) *UserHandler {
	return &UserHandler{repo: repo}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))

	var filters repository.UserFilter
	if idStr := q.Get("id"); idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		filters.ID = &id
	}
	if v := q.Get("name"); v != "" {
		filters.Name = &v
	}
	if v := q.Get("email"); v != "" {
		filters.Email = &v
	}
	if v := q.Get("gender"); v != "" {
		filters.Gender = &v
	}
	if v := q.Get("birth_date"); v != "" {
		filters.BirthDate = &v
	}
	filters.Status = q.Get("status")

	params := repository.UserListParams{
		Page:     page,
		PageSize: pageSize,
		OrderBy:  q.Get("order_by"),
		OrderDir: q.Get("order_dir"),
		Filters:  filters,
	}

	resp, err := h.repo.GetPaginatedUsers(params)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) GetUsersByCursor(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	limit, _ := strconv.Atoi(q.Get("limit"))

	var cursor *int
	if c := q.Get("cursor"); c != "" {
		id, err := strconv.Atoi(c)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid cursor")
			return
		}
		cursor = &id
	}

	params := repository.CursorParams{
		Limit:    limit,
		Cursor:   cursor,
		OrderBy:  q.Get("order_by"),
		OrderDir: q.Get("order_dir"),
		Status:   q.Get("status"),
	}

	resp, err := h.repo.GetUsersByCursor(params)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	user1ID, err := strconv.Atoi(q.Get("user1_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user1_id")
		return
	}

	user2ID, err := strconv.Atoi(q.Get("user2_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user2_id")
		return
	}

	users, err := h.repo.GetCommonFriends(user1ID, user2ID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user1_id":       user1ID,
		"user2_id":       user2ID,
		"common_friends": users,
	})
}

func (h *UserHandler) SoftDeleteUser(w http.ResponseWriter, r *http.Request) {
	// /users/{id}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 {
		writeError(w, http.StatusBadRequest, "invalid path")
		return
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	err = h.repo.SoftDeleteUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "user not found or already deleted")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "user soft deleted",
	})
}
