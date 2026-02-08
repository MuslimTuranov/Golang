package main

import (
	"Assignment1/internal/handlers"
	"Assignment1/internal/middleware"
	"Assignment1/internal/store"
	"log"
	"net/http"
)

func main() {
	storage := store.NewStorage()
	taskHandler := handlers.NewTaskHandler(storage)

	mux := http.NewServeMux()

	mux.Handle("/tasks", taskHandler)

	mux.HandleFunc("/external/get", taskHandler.HandleExternalGet)
	mux.HandleFunc("/external/post", taskHandler.HandleExternalPost)

	handlerWithAuth := middleware.Auth(mux)
	handlerWithReqID := middleware.RequestID(handlerWithAuth)
	finalHandler := middleware.Logging(handlerWithReqID)

	log.Println("starting at 8080 port")
	if err := http.ListenAndServe(":8080", finalHandler); err != nil {
		log.Fatal(err)
	}
}
