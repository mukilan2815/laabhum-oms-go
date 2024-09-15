package adapter

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Mukilan-T/laabhum-broker-adapter-go/internal/broker"
	"github.com/Mukilan-T/laabhum-broker-adapter-go/internal/cache"
	"github.com/Mukilan-T/laabhum-broker-adapter-go/internal/config"
	"github.com/Mukilan-T/laabhum-broker-adapter-go/internal/metrics"
	"github.com/Mukilan-T/laabhum-broker-adapter-go/pkg/sdk"
	"github.com/Mukilan-T/laabhum-broker-adapter-go/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
)

type Adapter struct {
	cfg               *config.Config
	brokerClient      *broker.Client
	kafkaProducer     *kafka.Writer
	kafkaConsumer     *kafka.Reader
	orderHandler      *OrderHandler
	posHandler        *PositionHandler
	marketDataHandler *MarketDataHandler
	cache             cache.Cache
	metrics           *metrics.Registry
	circuitBreaker    *utils.CircuitBreaker
}
func New(cfg *config.Config, metricsRegistry *metrics.Registry) (*Adapter, error) {
	// Create broker client
	brokerClient, err := broker.NewClient(&cfg.BrokerConfig)
	if err != nil {
		return nil, err
	}

	// Initialize Kafka Producer
	kafkaProducer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.KafkaConfig.Brokers,
		Topic:   cfg.KafkaConfig.Topic,
	})
	if err != nil {
		return nil, err
	}

	kafkaConsumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.KafkaConfig.Brokers,
		Topic:   cfg.KafkaConfig.Topic,
		GroupID: cfg.KafkaConfig.GroupID,
	})

	// Initialize Redis Cache
	redisCache, err := cache.NewRedisCache(&cfg.RedisConfig) // Make sure this matches the expected type
	if err != nil {
		return nil, err
	}

	// Create the adapter object
	adapter := &Adapter{
		cfg:            cfg,
		brokerClient:   brokerClient,
		kafkaProducer:  kafkaProducer,
		kafkaConsumer:  kafkaConsumer,
		cache:          redisCache,
		metrics:        metricsRegistry,
		circuitBreaker: utils.NewCircuitBreaker(cfg.CircuitBreakerConfig.MaxFailures, time.Duration(cfg.CircuitBreakerConfig.Timeout)*time.Second),
	}

	// Create specific handlers
	adapter.orderHandler = NewOrderHandler(*adapter.brokerClient, adapter.kafkaProducer)
	adapter.posHandler = NewPositionHandler(cfg.KafkaConfig.Brokers, cfg.KafkaConfig.Topic)
	adapter.marketDataHandler = NewMarketDataHandler()

	return adapter, nil
}

// Start starts the adapter services.
func (a *Adapter) Start() error {
	var wg sync.WaitGroup
	wg.Add(3)

	// Start Kafka consumer
	go func() {
		defer wg.Done()
		if err := a.startKafkaConsumer(); err != nil {
			log.Printf("Error starting Kafka consumer: %v", err)
		}
	}()

	// Start WebSocket connection
	go func() {
		defer wg.Done()
		if err := a.brokerClient.StartWebSocket(); err != nil {
			log.Printf("Error starting WebSocket: %v", err)
		}
	}()

	// Start periodic tasks (e.g. sync positions)
	go func() {
		defer wg.Done()
		a.startPeriodicTasks()
	}()

	wg.Wait()
	return nil
}

// Stop gracefully shuts down the adapter services.
func (a *Adapter) Stop(ctx context.Context) error {
	log.Println("Stopping adapter services...")

	if err := a.kafkaConsumer.Close(); err != nil {
		log.Printf("Error closing Kafka consumer: %v", err)
	}

	if err := a.brokerClient.Close(); err != nil {
		log.Printf("Error closing broker connection: %v", err)
	}

	return nil
}

// startKafkaConsumer starts the Kafka consumer for processing incoming messages.
func (a *Adapter) startKafkaConsumer() error {
	for {
		msg, err := a.kafkaConsumer.FetchMessage(context.Background())
		if err != nil {
			return err
		}

		// Process incoming Kafka messages
		log.Printf("Received Kafka message: %s", string(msg.Value))
		// Implement message processing logic here

		if err := a.kafkaConsumer.CommitMessages(context.Background(), msg); err != nil {
			return err
		}
	}
}

// startPeriodicTasks runs tasks like syncing positions periodically.
func (a *Adapter) startPeriodicTasks() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.syncPositions()
			a.checkOrderStatus()
		}
	}
}

// syncPositions syncs positions periodically.
func (a *Adapter) syncPositions() {
	log.Println("Syncing positions...")
	// Fetch positions from broker client
	positions, err := a.brokerClient.GetPositions()
	if err != nil {
		log.Printf("Error syncing positions: %v", err)
		return
	}

	// Cache positions in Redis
	if err := a.cache.Set("positions", positions, 10*time.Minute); err != nil {
		log.Printf("Error caching positions: %v", err)
	}
}

// checkOrderStatus checks the order status periodically.
func (a *Adapter) checkOrderStatus() {
	log.Println("Checking order status...")
	// Fetch order statuses from broker client
	orderStatuses, err := a.brokerClient.GetOrderStatuses()
	if err != nil {
		log.Printf("Error checking order status: %v", err)
		return
	}
	log.Printf("Order statuses: %v", orderStatuses)
}

// SetupRoutes sets up the HTTP routes for the adapter.
func (a *Adapter) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/order", a.CreateOrder).Methods("POST")
	router.HandleFunc("/positions", a.GetPositions).Methods("GET")
	router.HandleFunc("/marketdata", a.StreamMarketData).Methods("GET")
}

// CreateOrder handles the creation of an order.
func (a *Adapter) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order sdk.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	// Use OrderHandler to create the order
	resp, err := a.orderHandler.Create(r.Context(), order)
	if err != nil {
		http.Error(w, "failed to place order", http.StatusInternalServerError)
		return
	}

	// Respond with order confirmation
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// GetPositions handles fetching positions.
func (a *Adapter) GetPositions(w http.ResponseWriter, r *http.Request) {
	// Fetch positions from cache
	positions, err := a.cache.Get("positions")
	if err != nil {
		http.Error(w, "failed to fetch positions", http.StatusInternalServerError)
		return
	}

	// Respond with positions data
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(positions); err != nil {
		http.Error(w, "failed to encode positions", http.StatusInternalServerError)
	}
}

// StreamMarketData handles real-time market data streaming via WebSocket.
func (a *Adapter) StreamMarketData(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "failed to upgrade to websocket", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// Subscribe to market data
	marketDataCh, err := a.brokerClient.Subscribe("AAPL")
	if err != nil {
		http.Error(w, "failed to subscribe to market data", http.StatusInternalServerError)
		return
	}

	// Stream market data to WebSocket client
	for marketData := range marketDataCh {
		if err := conn.WriteJSON(marketData); err != nil {
			log.Printf("Error writing market data to websocket: %v", err)
			break
		}
	}
}
