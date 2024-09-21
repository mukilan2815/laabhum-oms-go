package oms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Order represents an order in the system
type Order struct {
	ID          string  `json:"id"`
	Symbol      string  `json:"symbol"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Side        string  `json:"side"`   // "buy" or "sell"
	Status      string  `json:"status"`
	CreatedAt   int64   `json:"created_at"` // Optional, for tracking creation time
	Description string  `json:"description,omitempty"` // Optional, use omitempty if not always needed
}

// Client is the OMS client structure
type Client struct {
	BaseURL string
}

type Position struct {
	ID           string  `json:"id"`           // Unique identifier for the position
	Symbol       string  `json:"symbol"`       // Trading symbol (e.g., stock, currency)
	Quantity     int     `json:"quantity"`     // Quantity of the asset in the position
	EntryPrice   float64 `json:"entry_price"`  // Price at which the position was entered
	CurrentPrice float64 `json:"current_price"`// Current market price of the asset
	Side         string  `json:"side"`         // Position side: "buy" or "sell"
	Status       string  `json:"status"`       // Status of the position (e.g., open, closed)
	Timestamp    int64   `json:"timestamp"`    // Timestamp of when the position was created
}

type PositionOrder struct {
	Symbol   string `json:"symbol"`
	Quantity int    `json:"quantity"`
}

// Convert PositionOrder to Order
func convertPositionOrderToOrder(po PositionOrder) Order {
	return Order{
		ID:          "", // Assign an appropriate ID as needed
		Description: po.Symbol,
		Quantity:    po.Quantity,
	}
}
func (c *Client) ExecuteChildOrder(parentID, childID string) ([]byte, error) {

    // Implement the logic to execute a child order

    // This is a placeholder implementation

    return []byte(`{"status": "success"}`), nil

}

// CreatePositionOrder creates a new position order
func (c *Client) CreatePositionOrder(order Order) error {
	url := fmt.Sprintf("%s/oms/positions/orders", c.BaseURL)
	body, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal position order: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create position order: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body) // Read response body for more details
		return fmt.Errorf("failed to create position order, status code: %d, body: %s", resp.StatusCode, body)
	}
	return nil
}

// GetOrders retrieves all orders
func (c *Client) GetOrders() ([]byte, error) {
	url := c.BaseURL + "/orders"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get orders, status code: %d, body: %s", resp.StatusCode, body)
	}

	return ioutil.ReadAll(resp.Body)
}

// Usage example function
// Handler represents a handler with an OMS client
type Handler struct {
	omsClient *Client
}

func ExampleUsage(h *Handler) {
	adjustedOrder := PositionOrder{
		Symbol:   "AAPL",
		Quantity: 10,
	}

	order := convertPositionOrderToOrder(adjustedOrder)
	err := h.omsClient.CreatePositionOrder(order)
	if err != nil {
		fmt.Println("Error creating position order:", err)
	} else {
		fmt.Println("Position order created successfully")
	}
}

func (c *Client) GetPositionsBySymbol(symbol string) ([]Position, error) {
	url := fmt.Sprintf("%s/oms/positions?symbol=%s", c.BaseURL, symbol)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions by symbol: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get positions by symbol, status code: %d", resp.StatusCode)
	}

	var positions []Position
	if err := json.NewDecoder(resp.Body).Decode(&positions); err != nil {
		return nil, fmt.Errorf("failed to decode positions response: %w", err)
	}

	return positions, nil
}

// NewClient creates a new OMS client
func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

func (c *Client) ExecuteOrder(orderID string) error {
	url := fmt.Sprintf("%s/oms/order/%s/execute", c.BaseURL, orderID)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to execute order: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to execute order, status code: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) CancelOrder(orderID string) error {
	url := fmt.Sprintf("%s/oms/order/%s/cancel", c.BaseURL, orderID)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to cancel order, status code: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) CreateOrder(order Order) ([]byte, error) {
    url := c.BaseURL + "/orders"
    body, err := json.Marshal(order)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal order: %w", err)
    }
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
    if err != nil {
        return nil, fmt.Errorf("failed to create order: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        return nil, fmt.Errorf("failed to create order, status code: %d, body: %s", resp.StatusCode, body)
    }

    return ioutil.ReadAll(resp.Body)
}


// Implement the other methods similarly...

// GetPositions retrieves current positions
func (c *Client) GetPositions() ([]byte, error) {
	url := c.BaseURL + "/oms/positions"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get positions, status code: %d, body: %s", resp.StatusCode, body)
	}

	return ioutil.ReadAll(resp.Body)
}

// SyncPositions syncs positions
func (c *Client) SyncPositions() error {
	url := c.BaseURL + "/oms/positions/sync"
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to sync positions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to sync positions, status code: %d", resp.StatusCode)
	}
	return nil
}

// ConvertPosition converts a position
func (c *Client) ConvertPosition(positionID string) error {
	url := fmt.Sprintf("%s/oms/positions/%s/convert", c.BaseURL, positionID)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to convert position: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to convert position, status code: %d", resp.StatusCode)
	}
	return nil
}

// DeletePositionOrder deletes a position order
func (c *Client) DeletePositionOrder(positionID string) error {
	url := fmt.Sprintf("%s/oms/positions/orders/%s", c.BaseURL, positionID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete position order: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete position order, status code: %d", resp.StatusCode)
	}
	return nil
}

// ExitAllTrades exits all trades
func (c *Client) ExitAllTrades() error {
	url := c.BaseURL + "/oms/scalper/exit/trade"
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to exit all trades: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to exit all trades, status code: %d", resp.StatusCode)
	}
	return nil
}

// CancelAllChildOrders cancels all child orders for a parent ID
func (c *Client) CancelAllChildOrders(parentID string) error {
	url := fmt.Sprintf("%s/oms/scalper/order/%s/cancel", c.BaseURL, parentID)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to cancel all child orders: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to cancel all child orders, status code: %d", resp.StatusCode)
	}
	return nil
}

// GetTrades retrieves trades based on parentID
func (c *Client) GetTrades(parentID string) ([]byte, error) {
	url := fmt.Sprintf("%s/oms/scalper/trades/%s", c.BaseURL, parentID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get trades: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get trades, status code: %d, body: %s", resp.StatusCode, body)
	}

	return ioutil.ReadAll(resp.Body)
}

// DeleteOrder deletes an order based on parentID
func (c *Client) DeleteOrder(parentID string) error {
	url := fmt.Sprintf("%s/oms/scalper/order/%s", c.BaseURL, parentID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete order, status code: %d", resp.StatusCode)
	}
	return nil
}
