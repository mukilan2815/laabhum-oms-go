package api

import (
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/laabhum/laabhum-oms-go/internal/service"
    "github.com/laabhum/laabhum-oms-go/pkg/kafka"
)

type Handlers struct {
    oms *service.OMS
}

func NewHandlers(oms *service.OMS) *Handlers {
    return &Handlers{oms: oms}
}

func RegisterHandlers(r *mux.Router, oms *service.OMS) {
    h := NewHandlers(oms)

    r.HandleFunc("/oms/scalper/order", h.CreateOrder).Methods("POST")
    r.HandleFunc("/oms/scalper/order/{id}/execute", h.ExecuteOrder).Methods("POST")
    r.HandleFunc("/oms/scalper/order/{id}/cancel", h.CancelOrder).Methods("POST")
}

func (h *Handlers) CreateOrder(w http.ResponseWriter, r *http.Request) {
    var orderRequest struct {
        Symbol   string  `json:"symbol"`
        Quantity int     `json:"quantity"`
        Price    float64 `json:"price"`
        Side     string  `json:"side"` // buy or sell
    }

    err := json.NewDecoder(r.Body).Decode(&orderRequest)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    order := h.oms.CreateOrder(orderRequest.Symbol, orderRequest.Quantity, orderRequest.Price, orderRequest.Side)
    kafka.PublishOrderCreated(order.ID)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(order)
}

func (h *Handlers) ExecuteOrder(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    orderID := vars["id"]

    order, success := h.oms.ExecuteOrder(orderID)
    if !success {
        http.Error(w, "Order not found or cannot be executed", http.StatusBadRequest)
        return
    }

    kafka.PublishOrderExecuted(order.ID)
    json.NewEncoder(w).Encode(order)
}

func (h *Handlers) CancelOrder(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    orderID := vars["id"]

    order, success := h.oms.CancelOrder(orderID)
    if !success {
        http.Error(w, "Order not found or cannot be canceled", http.StatusBadRequest)
        return
    }

    kafka.PublishOrderCanceled(order.ID)
    json.NewEncoder(w).Encode(order)
}
