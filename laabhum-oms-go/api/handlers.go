package api

import (
	"encoding/json"
	"net/http"

	"github.com/Mukilan-T/laabhum-oms-go/models"
	"github.com/Mukilan-T/laabhum-oms-go/repository"
	"github.com/Mukilan-T/laabhum-oms-go/service"
	"github.com/gorilla/mux"
)

// Handlers struct to hold OMSService
type Handlers struct {
	omsService *service.OMSService
}

// NewHandlers initializes the Handlers
func NewHandlers(omsService *service.OMSService) *Handlers {
	return &Handlers{
		omsService: omsService,
	}
}

// Helper function to bind JSON and handle errors
func bindJSON(w http.ResponseWriter, r *http.Request, obj interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(obj); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	return nil
}

// CreateOrder handles creating a new order
func (h *Handlers) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := bindJSON(w, r, &order); err != nil {
		return
	}

	createdOrder, err := h.omsService.CreateOrder(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdOrder)
}

// CreateScalperOrder handles the creation of a scalper order
func (h *Handlers) CreateScalperOrder(w http.ResponseWriter, r *http.Request) {
	var order models.ScalperOrder
	if err := bindJSON(w, r, &order); err != nil {
		return
	}

	createdOrder, err := h.omsService.CreateScalperOrder(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdOrder)
}

// ExecuteChildOrder handles executing a child order
func (h *Handlers) ExecuteChildOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	parentID := vars["parentID"]
	childID := vars["childID"]

	err := h.omsService.ExecuteChildOrder(parentID, childID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Child order executed successfully"})
}

// GetTrades handles fetching trades for a parent order
func (h *Handlers) GetTrades(w http.ResponseWriter, r *http.Request) {
	parentID := mux.Vars(r)["parentId"]

	trades, err := h.omsService.GetTrades(parentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trades)
}

// GetOrders handles fetching all orders
func (h *Handlers) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.omsService.GetOrders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// ModifyOrder handles modifying an order
func (h *Handlers) ModifyOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	parentID := vars["parentId"]
	childID := vars["childId"]
	var newData map[string]interface{}
	if err := bindJSON(w, r, &newData); err != nil {
		return
	}

	err := h.omsService.ModifyOrder(parentID, childID, newData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order modified successfully"})
}

// CancelOrder handles canceling an order
func (h *Handlers) CancelOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	parentID := vars["parentId"]
	orderID := vars["orderId"]

	err := h.omsService.CancelOrder(parentID, orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order canceled successfully"})
}

// SetupRoutes sets up the routes for the API
func SetupRoutes(repo repository.OrderRepository, omsService *service.OMSService) *mux.Router {
	router := mux.NewRouter()
	h := NewHandlers(omsService)

	// Order routes
	router.HandleFunc("/orders", h.CreateOrder).Methods(http.MethodPost)
	router.HandleFunc("/orders", h.GetOrders).Methods(http.MethodGet)

	// Scalper order routes
	router.HandleFunc("/oms/scalper/order", h.CreateScalperOrder).Methods(http.MethodPost)
	router.HandleFunc("/oms/scalper/order/{parentID}/{childID}/execute", h.ExecuteChildOrder).Methods(http.MethodPost)
	router.HandleFunc("/oms/scalper/trades/{parentId}", h.GetTrades).Methods(http.MethodGet)

	// Order modification routes
	router.HandleFunc("/oms/scalper/order/{parentId}/{childId}/modify", h.ModifyOrder).Methods(http.MethodPatch)
	router.HandleFunc("/oms/scalper/order/{parentId}/{orderId}/cancel", h.CancelOrder).Methods(http.MethodPost)

	return router
}
