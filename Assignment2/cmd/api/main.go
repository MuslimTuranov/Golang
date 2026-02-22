package main

import (
	_ "Assignment2/docs"
	api "Assignment2/internal"
	"Assignment2/internal/app"
	"Assignment2/internal/handlers"
	"Assignment2/internal/repository/postgres"
	userrepo "Assignment2/internal/repository/postgres/users"
	"Assignment2/internal/usecase"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Users API
// @version 1.0
// @description Assignment2
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-KEY
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	cfg := app.LoadConfig()

	ctx := context.Background()
	db := postgres.NewDialect(ctx, &cfg.DB)
	defer db.DB.Close()

	userRepo := userrepo.NewUserRepo(db)

	userUC := usecase.NewUserUsecase(userRepo)
	authUC := usecase.NewAuthUsecase(userRepo, cfg.JWTSecret, cfg.JWTTTL)
	h := handlers.NewHandler(userUC, authUC)

	router := api.NewRouter(h, cfg.APIKey, cfg.JWTSecret)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("server started on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("shutting down")
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
