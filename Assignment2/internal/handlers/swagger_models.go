package handlers

import "Assignment2/pkg/modules"

type healthResponse struct {
	Status string `json:"status" example:"ok"`
}

type errorResponse struct {
	Error string `json:"error" example:"invalid input"`
}

type idResponse struct {
	ID int `json:"id" example:"1"`
}

type statusResponse struct {
	Status string `json:"status" example:"updated"`
}

type deleteResponse struct {
	RowsAffected int64 `json:"rows_affected" example:"1"`
}

type tokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type usersListResponse []modules.User
