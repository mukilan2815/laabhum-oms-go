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

	"github.com/Mukilan-T/laabhum-gateway-go/api"
	"github.com/Mukilan-T/laabhum-gateway-go/config"
	"github.com/Mukilan-T/laabhum-gateway-go/internal/oms"
	"github.com/Mukilan-T/laabhum-gateway-go/pkg/logger"
)

type Order struct {
	ID          string  `json:"id"`
	Symbol      string  `json:"symbol"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Side        string  `json:"side"`   // "buy" or "sell"
	Status      string  `json:"status"`
	CreatedAt   int64   `json:"created_at"` // Optional, for tracking creation time
	Description string  `json:"description,omitempty"` // Optional, use omitempty if not always needed
}


func main() {
	cfg := config.LoadConfig()
	if cfg == nil {
		log.Fatalf("Failed to load configuration")
	}
	fmt.Printf("Loaded OMS address: %s\n", cfg.Oms.BaseURL)

	logLevel := cfg.LogLevel
	customLogger := logger.New(logLevel)

	stdLogger := log.New(customLogger.Writer(), "", log.LstdFlags)

	// Initialize OMS client
	omsClient := oms.NewClient(cfg.Oms.BaseURL)
	if omsClient == nil {
		stdLogger.Fatalf("Failed to create OMS client")
	}

	// Log orders 
	ordersData, err := omsClient.GetOrders()
	if err != nil {
		stdLogger.Fatalf("Failed to get orders: %v", err)
	}

	var orders []Order
	if err := json.Unmarshal(ordersData, &orders); err != nil {
		stdLogger.Fatalf("Failed to unmarshal orders: %v", err)
	}

	for _, order := range orders {
		stdLogger.Printf("Order ID: %s, Status: %s", order.ID, order.Status)
	}

	// Initialize strategy builder

	router := api.SetupRoutes(cfg, customLogger, omsClient)

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	go func() {
		stdLogger.Printf("Starting server on %s", cfg.ServerAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			stdLogger.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	stdLogger.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		stdLogger.Fatalf("Server forced to shutdown: %v", err)
	}

	stdLogger.Println("Server exiting")
}
