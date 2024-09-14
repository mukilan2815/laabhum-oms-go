package kafka

import (
    "log"
)

func PublishOrderCreated(orderID string) {
    log.Printf("Published 'Order Created' event for order ID: %s", orderID)
}

func PublishOrderExecuted(orderID string) {
    log.Printf("Published 'Order Executed' event for order ID: %s", orderID)
}

func PublishOrderCanceled(orderID string) {
    log.Printf("Published 'Order Canceled' event for order ID: %s", orderID)
}
