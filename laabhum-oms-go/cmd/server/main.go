package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Mukilan-T/laabhum-oms-go/api"
	"github.com/Mukilan-T/laabhum-oms-go/models"
	"github.com/Mukilan-T/laabhum-oms-go/repository"
	"github.com/Mukilan-T/laabhum-oms-go/service"
)

func main() {
	// Initialize repository and service
	repo := repository.NewInMemoryOrderRepository()
	omsService := service.NewOMSService(repo)

	// Set up routes
	router := api.SetupRoutes(repo, omsService)

	// Add global middleware
	router.Use(loggingMiddleware)
	router.Use(errorHandlingMiddleware)
	router.Use(metricsMiddleware)

	// Configure server
	srv := &http.Server{
		Addr:         ":8081",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			logError(err, "HTTP server Shutdown")
		}
	}()

	// Start the server
	logInfo("Server started at :8081")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

// Middleware for logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logInfo("Request completed", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))
	})
}

// Error handling middleware
func errorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logError(fmt.Errorf("%v", err), "Panic occurred")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Metrics middleware (example, you'd typically use a proper metrics library)
func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		logInfo("Request metrics", "duration", duration)
		// Here you would typically send this data to a metrics collection system
	})
}

// Helper function for structured logging
func logError(err error, message string) {
	log.Printf("ERROR: %s: %v", message, err)
}

// Helper function for structured info logging
func logInfo(message string, keyvals ...interface{}) {
	log.Printf("INFO: %s %v", message, keyvals)
}

// Helper function to send JSON response
func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logError(err, "Error encoding JSON response")
	}
}

// Validation functions
func validateScalperOrder(order models.ScalperOrder) error {
	if order.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if order.Quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}
	// Add more validation as needed
	return nil
}

func validateOrder(order models.Order) error {
	if order.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if order.Side != "buy" && order.Side != "sell" {
		return fmt.Errorf("side must be 'buy' or 'sell'")
	}
	if order.Quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}
	if order.Price <= 0 {
		return fmt.Errorf("price must be positive")
	}
	// Add more validation as needed
	return nil
}

func validatePositionOrder(order models.PositionOrder) error {
	if order.PositionID == "" {
		return fmt.Errorf("position ID is required")
	}
	if order.OrderType != "market" && order.OrderType != "limit" {
		return fmt.Errorf("order type must be 'market' or 'limit'")
	}
	if order.Quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}
	if order.OrderType == "limit" && order.Price <= 0 {
		return fmt.Errorf("price must be positive for limit orders")
	}
	// Add more validation as needed
	return nil
}

// Health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	sendJSONResponse(w, map[string]string{"status": "healthy"}, http.StatusOK)
}

// Basic authentication middleware
func basicAuthMiddleware(next http.HandlerFunc, username, password string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok || user != username || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}