package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mukilan-T/laabhum-gateway-go/config"
	"github.com/Mukilan-T/laabhum-gateway-go/internal/oms"
	"github.com/Mukilan-T/laabhum-gateway-go/pkg/logger"
	"github.com/Mukilan-T/laabhum-gateway-go/routes"
)

type Order struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	// Add other fields as necessary
}


func main() {
	// Load configuration
	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatalf("Failed to load configuration")
	}
	fmt.Printf("Loaded OMS address: %s\n", cfg.Oms.BaseURL)

	// Initialize logger
	logLevel := cfg.LogLevel
	logger := logger.New(logLevel)

	// Initialize OMS client
	omsClient := oms.NewClient(cfg.Oms.BaseURL)
	if omsClient == nil {
		logger.Fatalf("Failed to create OMS client")
	}

	// Log orders (example)
	ordersData, err := omsClient.GetOrders()
	if err != nil {
		logger.Fatalf("Failed to get orders: %v", err)
	}

	var orders []Order
	if err := json.Unmarshal(ordersData, &orders); err != nil {
		logger.Fatalf("Failed to unmarshal orders: %v", err)
	}

	for _, order := range orders {
		logger.Infof("Order ID: %s, Status: %s", order.ID, order.Status)
	}

	// Setup routes
	router := routes.SetupRoutes(cfg, logger, omsClient)

	// Create server
	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	// Start server
	go func() {
		logger.Infof("Starting server on %s", cfg.ServerAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exiting")
}
