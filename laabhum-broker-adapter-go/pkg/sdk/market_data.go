package sdk

import (
	"github.com/gorilla/websocket"
)

func SubscribeToMarketData(wsURL string) (*websocket.Conn, error) {
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return nil, err
	}
	return ws, nil
}
