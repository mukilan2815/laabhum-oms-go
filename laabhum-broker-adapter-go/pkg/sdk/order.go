package sdk

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type OrderType string
type OrderStatus string

const (
	OrderTypeMarket OrderType = "MARKET"
	OrderTypeLimit  OrderType = "LIMIT"
	// Add other order types

	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusCompleted OrderStatus = "COMPLETED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
	// Add other statuses
)

type Order struct {
	ID           string      `json:"id"`
	Symbol       string      `json:"symbol"`
	Quantity     int         `json:"quantity"`
	Price        float64     `json:"price"`
	Type         OrderType   `json:"type"`
	Side         string      `json:"side"`
	Status       OrderStatus `json:"status"`
	CreatedAt    time.Time   `json:"created_at"`
	LastModified time.Time   `json:"last_modified"`
}

type OrderResponse struct {
	OrderID     string `json:"order_id"` // Add this field
	Order       Order  `json:"order"`
	BrokerOrdID string `json:"broker_order_id"`
	Status      string `json:"status"` // Add this field to match the broker response
}

type OrderEventType string

const (
	OrderCreated   OrderEventType = "ORDER_CREATED"
	OrderUpdated   OrderEventType = "ORDER_UPDATED"
	OrderCancelled OrderEventType = "ORDER_CANCELLED"
)

type OrderEvent struct {
	Type     OrderEventType `json:"type"`
	Order    Order          `json:"order"`
	Response OrderResponse  `json:"response"`
}

func PlaceOrder(orderPayload map[string]interface{}, baseURL string) error {
	payload, _ := json.Marshal(orderPayload)
	req, err := http.NewRequest("POST", baseURL+"/order", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return err
	}

	_, err = http.DefaultClient.Do(req)
	return err
}

type MarketData struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Volume int     `json:"volume"`
}

func MapToSDKOrderResponse(brokerResponse map[string]interface{}) OrderResponse {
	return OrderResponse{
		OrderID:     brokerResponse["order_id"].(string), // Map the OrderID field
		BrokerOrdID: brokerResponse["broker_order_id"].(string),
		Status:      brokerResponse["status"].(string),
	}
}
