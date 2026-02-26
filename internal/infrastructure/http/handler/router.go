package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	appMiddleware "github.com/sample-provider/buy-credit-api/internal/infrastructure/http/middleware"
)

func SetupRouter(
	authHandler *AuthHandler,
	transactionHandler *TransactionHandler,
	authMiddleware *appMiddleware.AuthMiddleware,
) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API v1 routes
	r.Route("/v1", func(r chi.Router) {
		// Public routes
		r.Post("/auth/token", authHandler.CreateToken)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)

			// Transaction routes
			r.Post("/transactions", transactionHandler.CreateTransaction)
			r.Get("/transactions/{transactionId}", transactionHandler.GetTransaction)
		})
	})

	return r
}
