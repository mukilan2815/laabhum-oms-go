package service

import (
    "github.com/laabhum/laabhum-oms-go/internal/models"
    "github.com/laabhum/laabhum-oms-go/internal/repository"
    "time"
    "github.com/google/uuid"
)

type OMS struct {
    repo *repository.OrderRepository
}

func NewOMS(repo *repository.OrderRepository) *OMS {
    return &OMS{repo: repo}
}

func (oms *OMS) CreateOrder(symbol string, quantity int, price float64, side string) *models.Order {
    order := &models.Order{
        ID:        uuid.NewString(),
        Symbol:    symbol,
        Quantity:  quantity,
        Price:     price,
        Side:      side,
        Status:    "created",
        CreatedAt: time.Now().Unix(),
    }
    oms.repo.CreateOrder(order)
    return order
}

func (oms *OMS) ExecuteOrder(id string) (*models.Order, bool) {
    order, exists := oms.repo.GetOrder(id)
    if !exists || order.Status != "created" {
        return nil, false
    }
    order.Status = "executed"
    oms.repo.UpdateOrder(order)
    return order, true
}

func (oms *OMS) CancelOrder(id string) (*models.Order, bool) {
    order, exists := oms.repo.GetOrder(id)
    if !exists || order.Status != "created" {
        return nil, false
    }
    order.Status = "canceled"
    oms.repo.UpdateOrder(order)
    return order, true
}