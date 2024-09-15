package adapter

import (
	"net/http"

	"github.com/Mukilan-T/laabhum-broker-adapter-go/internal/broker"
	"github.com/Mukilan-T/laabhum-broker-adapter-go/pkg/utils"
)

type MarketDataHandler struct{}

func NewMarketDataHandler() *MarketDataHandler {
	return &MarketDataHandler{}
}

func (h *MarketDataHandler) StreamMarketData(w http.ResponseWriter, r *http.Request) {
	ws, err := broker.ConnectToMarketDataStream()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to connect to market data stream")
		return
	}

	// Forward market data stream to the client
	broker.StreamMarketData(ws, w, r)
}
