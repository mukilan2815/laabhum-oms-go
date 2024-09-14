package service

import (
	"time"

	"github.com/Mukilan-T/laabhum-oms-go/models"
	"github.com/Mukilan-T/laabhum-oms-go/repository"
	"github.com/google/uuid"
)

type OMSService struct {
    repo repository.OrderRepository
}

func NewOMSService(repo repository.OrderRepository) *OMSService {
    return &OMSService{
        repo: repo,
    }
}

func (s *OMSService) CreateScalperOrder(order models.ScalperOrder) (*models.ScalperOrder, error) {
    order.ID = uuid.NewString()
    order.CreatedAt = time.Now().Unix()
    return s.repo.CreateScalperOrder(order)
}

func (s *OMSService) ExecuteChildOrder(parentID, childID string) error {
    return s.repo.ExecuteChildOrder(parentID, childID)
}

func (s *OMSService) GetTrades(parentID string) ([]models.Trade, error) {
    return s.repo.GetTrades(parentID)
}

func (s *OMSService) CreateOrder(order models.Order) (*models.Order, error) {
    order.ID = uuid.NewString()
    order.CreatedAt = time.Now().Unix()
    return s.repo.CreateOrder(order)
}

func (s *OMSService) GetOrders() ([]models.Order, error) {
    return s.repo.GetOrders()
}