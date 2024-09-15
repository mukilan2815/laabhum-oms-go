package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Mukilan-T/laabhum-broker-adapter-go/internal/adapter"
	"github.com/Mukilan-T/laabhum-broker-adapter-go/internal/config"
	"github.com/Mukilan-T/laabhum-broker-adapter-go/internal/metrics"
	"github.com/gorilla/mux"
)
	
	func main() {
		// Load config
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}
	
		// Define Kafka broker address
		brokers := cfg.Kafka.Brokers
	
		// Prometheus metrics
		prometheus := metrics.NewPrometheus()
		prometheus.Setup()
	
		// Setup Router
		router := mux.NewRouter()
		// Remove or implement SetupRoutes if needed
		// adapter.SetupRoutes(router, kafkaProducer, prometheus)
	
		// Initialize position handler
		positionHandler := adapter.NewPositionHandler(brokers, "your-topic-name")  // Update to use correct function
	
		// Setup additional routes
		router.HandleFunc("/positions", positionHandler.GetPositions).Methods("GET")
		router.HandleFunc("/convert_position", positionHandler.ConvertPosition).Methods("POST")
	
		// Start server
		server := &http.Server{
			Handler:      router,
			Addr:         fmt.Sprintf(":%s", cfg.Server.Port),  // Ensure cfg.Server.Port is correctly set in config
			WriteTimeout: 15 * time.Second,
		}
	
		log.Printf("Starting server on port %s", cfg.Server.Port)
		log.Fatal(server.ListenAndServe())
	}
