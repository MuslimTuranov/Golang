package handlers

import "Assignment2/internal/usecase"

type Handler struct {
	Users *usecase.UserUsecase
	Auth  *usecase.AuthUsecase
}

func NewHandler(users *usecase.UserUsecase, auth *usecase.AuthUsecase) *Handler {
	return &Handler{Users: users, Auth: auth}
}
