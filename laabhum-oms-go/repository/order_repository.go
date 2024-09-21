package repository

import (
    "errors"
    "fmt"

    "github.com/Mukilan-T/laabhum-oms-go/models"
)

type OrderRepository interface {
    CreateOrder(order models.Order) (*models.Order, error)
    UpdateOrder(order *models.Order) error
    GetOrders() ([]models.Order, error)
    CreateScalperOrder(order models.ScalperOrder) (*models.ScalperOrder, error)
    ExecuteChildOrder(parentID, childID string) error
    GetTrades(parentID string) ([]models.Trade, error)
    GetOrder(id string) (*models.Order, error)
    SaveOrder(order *models.Order) error
}

type InMemoryOrderRepository struct {
    orders        map[string]*models.Order
    scalperOrders map[string]*models.ScalperOrder
    trades        map[string][]models.Trade
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
    return &InMemoryOrderRepository{
        orders:        make(map[string]*models.Order),
        scalperOrders: make(map[string]*models.ScalperOrder),
        trades:        make(map[string][]models.Trade),
    }
}

func (r *InMemoryOrderRepository) CreateOrder(order models.Order) (*models.Order, error) {
    r.orders[order.ID] = &order
    return &order, nil
}

func (r *InMemoryOrderRepository) CreateScalperOrder(order models.ScalperOrder) (*models.ScalperOrder, error) {
    r.scalperOrders[order.ID] = &order
    return &order, nil
}

func (r *InMemoryOrderRepository) ExecuteChildOrder(parentID, childID string) error {
    if parentOrder, exists := r.scalperOrders[parentID]; exists {
        if childOrder, exists := r.scalperOrders[childID]; exists {
            childOrder.Status = "executed"
            parentOrder.Status = "partially executed"
            return nil
        }
    }
    return errors.New("order not found")
}

func (r *InMemoryOrderRepository) GetTrades(parentID string) ([]models.Trade, error) {
    return r.trades[parentID], nil
}

func (r *InMemoryOrderRepository) GetOrders() ([]models.Order, error) {
    var orders []models.Order
    for _, order := range r.orders {
        orders = append(orders, *order)
    }
    return orders, nil
}

func (r *InMemoryOrderRepository) GetOrder(id string) (*models.Order, error) {
    order, exists := r.orders[id]
    if !exists {
        return nil, errors.New("order not found")
    }
    return order, nil
}

func (r *InMemoryOrderRepository) UpdateOrder(order *models.Order) error {
    r.orders[order.ID] = order
    return nil
}

func (r *InMemoryOrderRepository) SaveOrder(order *models.Order) error {
    if order == nil || order.ID == "" {
        return errors.New("invalid order")
    }
    r.orders[order.ID] = order
    return nil
}

// SaveOrder stores the order in the database
func SaveOrder(order map[string]interface{}) error {
    // This is just a stub. Replace with actual DB code
    fmt.Println("Order saved:", order)
    return nil
}