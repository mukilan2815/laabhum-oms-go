package repository

import (
    "github.com/laabhum/laabhum-oms-go/internal/models"
    "sync"
)

type OrderRepository struct {
    orders map[string]*models.Order
    mu     sync.Mutex
}

func NewOrderRepository() *OrderRepository {
    return &OrderRepository{
        orders: make(map[string]*models.Order),
    }
}

func (r *OrderRepository) CreateOrder(order *models.Order) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.orders[order.ID] = order
}

func (r *OrderRepository) GetOrder(id string) (*models.Order, bool) {
    r.mu.Lock()
    defer r.mu.Unlock()
    order, exists := r.orders[id]
    return order, exists
}

func (r *OrderRepository) UpdateOrder(order *models.Order) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.orders[order.ID] = order
}

func (r *OrderRepository) DeleteOrder(id string) {
    r.mu.Lock()
    defer r.mu.Unlock()
    delete(r.orders, id)
}
