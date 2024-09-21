package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"fmt"

	"github.com/Mukilan-T/laabhum-gateway-go/config"
	"github.com/Mukilan-T/laabhum-gateway-go/internal/oms"
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
func (h *Handlers) CreatePositionOrder(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Errorf("Failed to read request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var positionOrder oms.PositionOrder
	if err := json.Unmarshal(body, &positionOrder); err != nil {
		h.logger.Errorf("Failed to unmarshal position order: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid position order data"})
		return
	}

	if err := validatePositionOrder(positionOrder); err != nil {
		h.logger.Errorf("Position order validation failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingPositions, err := h.omsClient.GetPositionsBySymbol(positionOrder.Symbol)
	if err != nil {
		h.logger.Errorf("Failed to get existing positions: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get existing positions"})
		return
	}

	adjustedOrder := adjustPositionOrder(positionOrder, existingPositions)

	// Convert the adjusted position order to a standard order
	orderToCreate := oms.Order{
		ID:          "", // Assign a new ID as needed or let the backend generate it
		Description: adjustedOrder.Symbol,
		Quantity:    adjustedOrder.Quantity,
	}

	// Create the order using the omsClient
	err = h.omsClient.CreatePositionOrder(orderToCreate)
	if err != nil {
		h.logger.Errorf("Failed to create position order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create position order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Position order created successfully"})
}



// validatePositionOrder validates the position order
func validatePositionOrder(order oms.PositionOrder) error {
	if order.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if order.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}
	// Add more validation logic as needed
	return nil
}

// adjustPositionOrder adjusts the position order based on existing positions
func adjustPositionOrder(order oms.PositionOrder, existingPositions []oms.Position) oms.PositionOrder {
	// Example logic: Adjust the quantity based on existing positions
	for _, position := range existingPositions {
		if position.Symbol == order.Symbol {
			order.Quantity += position.Quantity
		}
	}
	return order
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
h.logger.Infof("Creating order: %+v", order)

	c.Data(http.StatusCreated, "application/json", response)
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

func (h *Handlers) GetOrders(c *gin.Context) {
	response, err := h.omsClient.GetOrders()
	if err != nil {
		h.logger.Errorf("Failed to get orders: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders"})
		return
	}

	c.Data(http.StatusOK, "application/json", response)
}

// SetupRoutes sets up the routes for the API
func SetupRoutes(cfg *config.Config, logger *logger.Logger, omsClient *oms.Client) *mux.Router {
	router := mux.NewRouter()
	h := NewHandlers(cfg, logger, omsClient)

	router.Handle("/orders", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := gin.CreateTestContext(w)
		c.Request = r
		h.CreateOrder(c)
	})).Methods(http.MethodPost)
	router.Handle("/scalper/orders", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := gin.CreateTestContext(w)
		c.Request = r
		h.CreateScalperOrder(c)
	})).Methods(http.MethodPost)
	router.Handle("/orders", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := gin.CreateTestContext(w)
		c.Request = r
		h.GetOrders(c)
	})).Methods(http.MethodGet)

	// Additional routes for child orders and trades
	router.Handle("/oms/child/{parentID}/{childID}/execute", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := gin.CreateTestContext(w)
		c.Request = r
		h.ExecuteChildOrder(c)
	})).Methods(http.MethodPost)

	return router
}
