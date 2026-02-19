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
	walletRepo := repository.NewInMemoryWalletRepository()
	transactionRepo := repository.NewInMemoryTransactionRepository()
	webhookRepo := repository.NewInMemoryWebhookRepository()
	partnerRepo := repository.NewInMemoryPartnerRepository()

	// Initialize JWT service
	jwtService := auth.NewJWTService("your-secret-key-change-in-production")

	// Initialize use cases
	authUseCase := application.NewAuthUseCase(partnerRepo, jwtService)
	walletUseCase := application.NewWalletUseCase(walletRepo)
	transactionUseCase := application.NewTransactionUseCase(transactionRepo, walletRepo, webhookRepo)
	webhookUseCase := application.NewWebhookUseCase(webhookRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUseCase)
	walletHandler := handler.NewWalletHandler(walletUseCase)
	transactionHandler := handler.NewTransactionHandler(transactionUseCase)
	webhookHandler := handler.NewWebhookHandler(webhookUseCase)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Setup router
	router := handler.SetupRouter(
		authHandler,
		walletHandler,
		transactionHandler,
		webhookHandler,
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
