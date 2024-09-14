package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Mukilan-T/laabhum-gateway-go/config"
	"github.com/Mukilan-T/laabhum-gateway-go/internal/oms"
	"github.com/Mukilan-T/laabhum-gateway-go/pkg/logger"
	"github.com/gin-gonic/gin"
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

	response, err := h.omsClient.CreateScalperOrder(json.RawMessage(body))
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
	parentID := c.Param("parentId")

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

	response, err := h.omsClient.CreateOrder(json.RawMessage(body))
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