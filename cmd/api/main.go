package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sample-provider/buy-credit-api/internal/application"
	"github.com/sample-provider/buy-credit-api/internal/infrastructure/auth"
	"github.com/sample-provider/buy-credit-api/internal/infrastructure/http/handler"
	"github.com/sample-provider/buy-credit-api/internal/infrastructure/http/middleware"
	"github.com/sample-provider/buy-credit-api/internal/infrastructure/repository"
)

func main() {
	// Initialize repositories (in-memory for this example)
	transactionRepo := repository.NewInMemoryTransactionRepository()
	partnerRepo := repository.NewInMemoryPartnerRepository()

	// Initialize JWT service
	jwtService := auth.NewJWTService("your-secret-key-change-in-production")

	// Initialize use cases
	authUseCase := application.NewAuthUseCase(partnerRepo, jwtService)
	transactionUseCase := application.NewTransactionUseCase(transactionRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUseCase)
	transactionHandler := handler.NewTransactionHandler(transactionUseCase)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Setup router
	router := handler.SetupRouter(
		authHandler,
		transactionHandler,
		authMiddleware,
	)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
