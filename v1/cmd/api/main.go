package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"dailystep-backend/internal/config"
	"dailystep-backend/internal/db"
	"dailystep-backend/internal/handler"
	appmiddleware "dailystep-backend/internal/middleware"
	"dailystep-backend/internal/repository"
	"dailystep-backend/internal/service"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	pool, err := db.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db error: %v", err)
	}
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, time.Duration(cfg.TokenTTLHours)*time.Hour)

	authHandler := handler.NewAuthHandler(authService)
	healthHandler := handler.NewHealthHandler()

	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(30 * time.Second))

	r.Get("/health", healthHandler.Check)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	r.Group(func(r chi.Router) {
		r.Use(appmiddleware.AuthMiddleware(cfg.JWTSecret))

		r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			userID, _ := appmiddleware.GetUserID(r.Context())
			handler.WriteMePlaceholder(w, userID)
		})
	})

	addr := ":" + cfg.HTTPPort
	log.Printf("server listening on %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
