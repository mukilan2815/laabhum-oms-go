package broker

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/Mukilan-T/laabhum-broker-adapter-go/internal/config"
	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client with subscription management.
type Client struct {
	wsClient      *websocket.Conn
	brokerClient  *BrokerClient
	subscriptions map[string]chan MarketData
	mu            sync.Mutex
}

// MarketData represents market data received from the WebSocket server.
type MarketData struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Volume int     `json:"volume"`
}

// BrokerClient represents the client for interacting with the broker's HTTP API.
type BrokerClient struct {
	BaseURL string
	APIKey  string
}

// Position represents a position held by the client.
type Position struct {
	// Define fields for Position
}

// OrderStatus represents the status of an order.
type OrderStatus struct {
	// Define fields for OrderStatus
}

// NewBrokerClient creates a new BrokerClient instance.
func NewBrokerClient(baseURL, apiKey string) *BrokerClient {
	return &BrokerClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
	}
}

// PlaceOrder sends an order to the broker's API.
func (client *BrokerClient) PlaceOrder(orderPayload map[string]interface{}) (*http.Response, error) {
	payload, _ := json.Marshal(orderPayload)
	req, err := http.NewRequest("POST", client.BaseURL+"/orders/place", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+client.APIKey)
	req.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}

// NewClient initializes a new WebSocket client with broker client.
func NewClient(cfg *config.BrokerConfig) (*Client, error) {
	wsClient, _, err := websocket.DefaultDialer.Dial(cfg.WebSocketURL, nil)
	if err != nil {
		return nil, err
	}

	brokerClient := NewBrokerClient(cfg.APIBaseURL, cfg.APIKey)

	return &Client{
		wsClient:      wsClient,
		brokerClient:  brokerClient,
		subscriptions: make(map[string]chan MarketData),
	}, nil
}

// StartWebSocket starts the WebSocket connection.
func (c *Client) StartWebSocket() error {
	go func() {
		for {
			_, message, err := c.wsClient.ReadMessage()
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}

			var marketData MarketData
			if err := json.Unmarshal(message, &marketData); err != nil {
				log.Println("Error unmarshalling message:", err)
				continue
			}

			c.handleMarketData(marketData)
		}
	}()
	return nil
}

// Close closes the WebSocket connection.
func (c *Client) Close() error {
	return c.wsClient.Close()
}

// Subscribe subscribes to market data for a specific symbol.
func (c *Client) Subscribe(symbol string) (chan MarketData, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.subscriptions[symbol]; exists {
		return nil, errors.New("already subscribed to symbol")
	}

	ch := make(chan MarketData, 100)
	c.subscriptions[symbol] = ch

	err := c.sendSubscribeRequest(symbol)
	if err != nil {
		delete(c.subscriptions, symbol)
		return nil, err
	}

	return ch, nil
}

// Unsubscribe unsubscribes from market data for a specific symbol.
func (c *Client) Unsubscribe(symbol string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.subscriptions[symbol]; !exists {
		return errors.New("not subscribed to symbol")
	}

	err := c.sendUnsubscribeRequest(symbol)
	if err != nil {
		return err
	}

	close(c.subscriptions[symbol])
	delete(c.subscriptions, symbol)

	return nil
}

// PlaceOrder places an order using the broker client.
func (c *Client) PlaceOrder(orderPayload map[string]interface{}) (*http.Response, error) {
	return c.brokerClient.PlaceOrder(orderPayload)
}

// GetPositions gets positions from the broker.
func (c *Client) GetPositions() ([]Position, error) {
	// Implement the logic to get positions from the broker
	return nil, nil
}

// GetOrderStatuses gets order statuses from the broker.
func (c *Client) GetOrderStatuses() ([]OrderStatus, error) {
	// Implement the logic to get order statuses from the broker
	return nil, nil
}

// sendSubscribeRequest sends a subscription request to the WebSocket server.
func (c *Client) sendSubscribeRequest(symbol string) error {
	message := map[string]string{
		"action": "subscribe",
		"symbol": symbol,
	}
	return c.wsClient.WriteJSON(message)
}

// sendUnsubscribeRequest sends an unsubscription request to the WebSocket server.
func (c *Client) sendUnsubscribeRequest(symbol string) error {
	message := map[string]string{
		"action": "unsubscribe",
		"symbol": symbol,
	}
	return c.wsClient.WriteJSON(message)
}

// handleMarketData processes incoming market data and routes it to the appropriate subscription channels.
func (c *Client) handleMarketData(marketData MarketData) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ch, exists := c.subscriptions[marketData.Symbol]; exists {
		select {
		case ch <- marketData:
		default:
			log.Printf("Market data channel for %s is full, dropping data", marketData.Symbol)
		}
	}
}
