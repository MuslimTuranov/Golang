package main

import (
	"log"
	"net/http"
	"users-service/internal/db"
	"users-service/internal/handler"
	"users-service/internal/repository"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	pg, err := db.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository(pg)
	h := handler.NewUserHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", h.GetUsers)
	mux.HandleFunc("GET /users/cursor", h.GetUsersByCursor)
	mux.HandleFunc("GET /users/common-friends", h.GetCommonFriends)
	mux.HandleFunc("DELETE /users/", h.SoftDeleteUser)

	log.Println("server running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
