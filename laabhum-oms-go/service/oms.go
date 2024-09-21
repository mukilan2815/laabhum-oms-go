package service

import (
	"errors"
	"time"

	"github.com/Mukilan-T/laabhum-oms-go/models"
	"github.com/Mukilan-T/laabhum-oms-go/repository"
	"github.com/google/uuid"
)

type OMSService struct {
	repo repository.OrderRepository
}

func NewOMSService(repo repository.OrderRepository) *OMSService {
	return &OMSService{repo: repo}
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

// ProcessOrder handles business logic for processing the order
func (s *OMSService) ProcessOrder(order map[string]interface{}) error {
	// Add business logic for order processing here
	if len(order) == 0 {
		return errors.New("invalid order data")
	}

	// Save the order to the repository (database)
	

	return nil
}
func (s *OMSService) PlaceOrder(order *models.Order) error {

    // Implement the logic to place an order

    return nil

}
// func (r *InMemoryOrderRepository) GetOrder(id string) (*models.Order, error) {

//     order, exists := r.orders[id]

//     if !exists {

//         return nil, fmt.Errorf("order not found")

//     }

//     return order, nil

// }

func (s *OMSService) ModifyOrder(parentID, childID string, newData map[string]interface{}) error {
	// Fetch the existing order
	order, err := s.repo.GetOrder(parentID)
	if err != nil {
		return err
	}

	// Update the order with new data
	for key, value := range newData {
		switch key {
		case "status":
			order.Status = value.(string)
		case "quantity":
			order.Quantity = value.(int)
		case "price":
			order.Price = value.(float64)
		// Add more fields as needed
		default:
			return errors.New("invalid field in newData")
		}
	}

	// Save the updated order back to the repository
	err = s.repo.UpdateOrder(order)
	if err != nil {
		return errors.New("failed to update order")
	}

	return nil
}

func (s *OMSService) CancelOrder(parentID, orderID string) error {
	// Fetch the existing order
	order, err := s.repo.GetOrder(parentID)
	if err != nil {
		return err
	}

	// Check if the order can be canceled
	if order.Status == "executed" {
		return errors.New("cannot cancel an executed order")
	}

	// Update the order status to canceled
	order.Status = "canceled"

	// Save the updated order back to the repository
	err = s.repo.UpdateOrder(order)
	if err != nil {
		return errors.New("failed to update order")
	}

	return nil
}