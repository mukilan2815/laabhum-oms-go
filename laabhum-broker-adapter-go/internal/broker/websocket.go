package broker

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Removed duplicate MarketData struct declaration

type WebSocketClient struct {
	conn         *websocket.Conn
	sendCh       chan []byte
	onMarketData func(MarketData)
}

func NewWebSocketClient(url string, onMarketData func(MarketData)) (*WebSocketClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	client := &WebSocketClient{
		conn:         conn,
		sendCh:       make(chan []byte, 100),
		onMarketData: onMarketData,
	}

	go client.readPump()
	go client.writePump()

	return client, nil
}

func (c *WebSocketClient) readPump() {
	defer c.conn.Close()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			return
		}

		var marketData MarketData
		if err := json.Unmarshal(message, &marketData); err != nil {
			log.Printf("Error unmarshalling market data: %v", err)
			continue
		}

		c.onMarketData(marketData)
	}
}

func (c *WebSocketClient) writePump() {
	ticker := time.NewTicker(time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.sendCh:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("WebSocket ping error: %v", err)
				return
			}
		}
	}
}

func (c *WebSocketClient) Send(ctx context.Context, message []byte) error {
	select {
	case c.sendCh <- message:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *WebSocketClient) Close() error {
	close(c.sendCh)
	return c.conn.Close()
}

// New function to connect to market data stream
func ConnectToMarketDataStream() (*websocket.Conn, error) {
	url := "wss://broker-stream-url.com/marketdata"
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

// New function to stream market data to an HTTP client
func StreamMarketData(ws *websocket.Conn, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	clientWS, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to upgrade: %v", err)
		return
	}
	defer clientWS.Close()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("error reading websocket message: %v", err)
			return
		}
		clientWS.WriteMessage(websocket.TextMessage, msg)
	}
}
