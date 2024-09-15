package adapter

import (
	"encoding/json"
	"log"
	"net/http"

	kafka "github.com/Mukilan-T/laabhum-broker-adapter-go/pkg/kafka/producer"
	"github.com/Mukilan-T/laabhum-broker-adapter-go/pkg/utils"
)

type PositionHandler struct {
    kafkaProducer *kafka.KafkaProducer
}

func NewPositionHandler(brokers []string, topic string) *PositionHandler {
    producer := kafka.NewProducer(brokers, topic)
    if producer == nil {
        // Handle error, e.g., log and exit or return nil
        log.Fatalf("Failed to create Kafka producer")
    }
    return &PositionHandler{kafkaProducer: producer}
}

func (h *PositionHandler) GetPositions(w http.ResponseWriter, r *http.Request) {
    positions := map[string]interface{}{
        "positions": []map[string]interface{}{{"symbol": "AAPL", "qty": 10}},
    }
    utils.RespondWithJSON(w, http.StatusOK, positions)
}

func (h *PositionHandler) ConvertPosition(w http.ResponseWriter, r *http.Request) {
    var conversionRequest map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&conversionRequest); err != nil {
        utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    utils.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "position converted"})
}
