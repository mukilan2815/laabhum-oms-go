package adapter

import (
	"context"
	"encoding/json"

	"github.com/Mukilan-T/laabhum-broker-adapter-go/internal/broker"
	"github.com/Mukilan-T/laabhum-broker-adapter-go/pkg/sdk"
	"github.com/segmentio/kafka-go"
)

type OrderHandler struct {
	brokerClient broker.Client
	kafkaProducer *kafka.Writer
}

func NewOrderHandler(brokerClient broker.Client, kafkaProducer *kafka.Writer) *OrderHandler {
	return &OrderHandler{brokerClient: brokerClient, kafkaProducer: kafkaProducer}
}

func (h *OrderHandler) Create(ctx context.Context, order sdk.Order) (sdk.OrderResponse, error) {
	// Convert SDK order to broker-specific order
	brokerOrder := broker.MapToBrokerOrder(order)

	// Send order to broker
	resp, err := h.brokerClient.PlaceOrder(brokerOrder)
	if err != nil {
		return sdk.OrderResponse{}, err
	}

	// Convert broker response to SDK response
	brokerResp := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&brokerResp); err != nil {
		return sdk.OrderResponse{}, err
	}
	sdkResp := broker.MapToSDKOrderResponse(brokerResp)

	// Publish order creation event to Kafka
	event := sdk.OrderEvent{
		Type:     sdk.OrderCreated,
		Order:    order,
		Response: sdkResp,
	}
	if err := h.kafkaProducer.WriteMessages(ctx, kafka.Message{
		Value: func() []byte {
			value, err := json.Marshal(event)
			if err != nil {
				// Log error, but don't fail the request
				// Consider implementing a retry mechanism
				return nil
			}
			return value
		}(),
	}); err != nil {
		// Log error, but don't fail the request
		// Consider implementing a retry mechanism
	}

	return sdkResp, nil
}
