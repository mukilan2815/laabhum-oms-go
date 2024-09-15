package broker

import (
	"github.com/Mukilan-T/laabhum-broker-adapter-go/pkg/sdk"
)

func MapToBrokerOrder(order sdk.Order) map[string]interface{} {
	return map[string]interface{}{
		"symbol": order.Symbol,
		"qty":    order.Quantity,
		"price":  order.Price,
		"type":   order.Type,
	}
}

func MapToSDKOrderResponse(brokerResponse map[string]interface{}) sdk.OrderResponse {
	return sdk.OrderResponse{
		OrderID: brokerResponse["order_id"].(string),
		Status:  brokerResponse["status"].(string),
	}
}

func MapOrderResponse(brokerResponse map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"order_id": brokerResponse["orderId"],
		"status":   brokerResponse["status"],
	}
}

func MapPositionResponse(brokerResponse map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"symbol": brokerResponse["symbol"],
		"qty":    brokerResponse["qty"],
		"price":  brokerResponse["price"],
	}
}
