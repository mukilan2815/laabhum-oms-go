package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Mukilan-T/laabhum-gateway-go/config"
	"github.com/Mukilan-T/laabhum-gateway-go/internal/oms"
	"github.com/Mukilan-T/laabhum-gateway-go/internal/strategy"
	"github.com/Mukilan-T/laabhum-gateway-go/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
)

type Handlers struct {
	cfg       *config.Config
	logger    *logger.Logger
	omsClient *oms.Client
}

func NewHandlers(cfg *config.Config, logger *logger.Logger, omsClient *oms.Client) *Handlers {
	return &Handlers{
		cfg:       cfg,
		logger:    logger,
		omsClient: omsClient,
	}
}

func (h *Handlers) CreateScalperOrder(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Errorf("Failed to read request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var order oms.Order
	if err := json.Unmarshal(body, &order); err != nil {
		h.logger.Errorf("Failed to unmarshal order: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order data"})
		return
	}

	response, err := h.omsClient.CreateOrder(order)
	if err != nil {
		h.logger.Errorf("Failed to create scalper order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create scalper order"})
		return
	}

	c.Data(http.StatusCreated, "application/json", response)
}

func (h *Handlers) ExecuteChildOrder(c *gin.Context) {
	parentID := c.Param("parentID")
	childID := c.Param("childID")

	response, err := h.omsClient.ExecuteChildOrder(parentID, childID)
	if err != nil {
		h.logger.Errorf("Failed to execute child order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute child order"})
		return
	}

	c.Data(http.StatusOK, "application/json", response)
}

func (h *Handlers) GetTrades(c *gin.Context) {
	parentID := c.Param("parentID")

	response, err := h.omsClient.GetTrades(parentID)
	if err != nil {
		h.logger.Errorf("Failed to get trades: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trades"})
		return
	}

	c.Data(http.StatusOK, "application/json", response)
}

func (h *Handlers) CreateOrder(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Errorf("Failed to read request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var order oms.Order
	if err := json.Unmarshal(body, &order); err != nil {
		h.logger.Errorf("Failed to unmarshal order: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order data"})
		return
	}

	response, err := h.omsClient.CreateOrder(order)
	if err != nil {
		h.logger.Errorf("Failed to create order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.Data(http.StatusCreated, "application/json", response)
}

func (h *Handlers) GetOrders(c *gin.Context) {
	response, err := h.omsClient.GetOrders()
	if err != nil {
		h.logger.Errorf("Failed to get orders: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}

	c.Data(http.StatusOK, "application/json", response)
}

// Handler for creating an order
func CreateOrderHandler(cfg *config.Config, omsClient *oms.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var orderRequest map[string]interface{}

		err := json.NewDecoder(r.Body).Decode(&orderRequest)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Call OMS client to create an order
		order := oms.Order{
			// Populate the order fields from orderRequest
			// Example: ID: orderRequest["id"].(string)
		}
		_, err = omsClient.CreateOrder(order)
		if err != nil {
			http.Error(w, "Failed to create order", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Order created successfully"})
	}
}

func SetupRoutes(cfg *config.Config, logger *logger.Logger, omsClient *oms.Client, strategyBuilder *strategy.Builder) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			ordersData, err := omsClient.GetOrders()
			if err != nil {
				logger.Errorf("Failed to get orders: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(ordersData)
			return
		}
		if r.Method == http.MethodPost {
			var order oms.Order
			if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
				logger.Errorf("Failed to decode order: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			strategyResponse, err := strategyBuilder.ProcessOrder(order)
			if err != nil {
				logger.Errorf("Failed to process order: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if strategyResponse != "" {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(strategyResponse))
				return
			}
			createdOrder, err := omsClient.CreateOrder(order)
			if err != nil {
				logger.Errorf("Failed to create order: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(createdOrder)
		}
	}).Methods(http.MethodGet, http.MethodPost)
	// Add the new route for CreateOrderHandler
	router.HandleFunc("/create-order", CreateOrderHandler(cfg, omsClient)).Methods(http.MethodPost)
	router.HandleFunc("/create-order", CreateOrderHandler(cfg, omsClient)).Methods(http.MethodPost)

	return router
}
