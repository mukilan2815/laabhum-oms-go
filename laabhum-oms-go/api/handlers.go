package api

import (
	"net/http"

	"github.com/Mukilan-T/laabhum-oms-go/models"
	"github.com/Mukilan-T/laabhum-oms-go/service"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	omsService *service.OMSService
}

func NewHandlers(omsService *service.OMSService) *Handlers {
	return &Handlers{
		omsService: omsService,
	}
}

func (h *Handlers) CreateScalperOrder(c *gin.Context) {
	var order models.ScalperOrder
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdOrder, err := h.omsService.CreateScalperOrder(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdOrder)
}

func (h *Handlers) ExecuteChildOrder(c *gin.Context) {
	parentID := c.Param("parentID")
	childID := c.Param("childID")

	err := h.omsService.ExecuteChildOrder(parentID, childID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Child order executed successfully"})
}

func (h *Handlers) GetTrades(c *gin.Context) {
	parentID := c.Param("parentId")

	trades, err := h.omsService.GetTrades(parentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trades)
}

func (h *Handlers) CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdOrder, err := h.omsService.CreateOrder(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdOrder)
		c.JSON(http.StatusCreated, gin.H{"message": "Order created successfully"})

}

func (h *Handlers) GetOrders(c *gin.Context) {
	orders, err := h.omsService.GetOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}