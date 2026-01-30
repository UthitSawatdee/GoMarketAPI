package usecase

import (
	// "errors"
	// "fmt"
	// domain "github.com/Fal2o/E-Commerce_API/internal/domain"
	"errors"

	"github.com/Fal2o/E-Commerce_API/internal/domain"
	port "github.com/Fal2o/E-Commerce_API/internal/port"
	// "gorm.io/gorm"
)

// OrderUseCase defines the interface for user business logic
// คุยกับ service (fiber)
type OrderUseCase interface {
	ViewOrder(userID uint) ([]*domain.Order, error)
	CancelOrder(orderID string) error
	AllOrders() ([]*domain.Order, error)
	UpdateOrderStatus(orderID string, status string) (*domain.Order, string, error)
}

type OrderService struct {
	repo port.OrderRepository
}

func NewOrderService(repo port.OrderRepository) OrderUseCase {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) ViewOrder(userID uint) ([]*domain.Order, error) {
	// 1. Check if email already exists
	order, err := s.repo.GetOrderByUserID(userID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) CancelOrder(orderID string) error {
	err := s.repo.DeleteOrderByOrderID(orderID)
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderService) AllOrders() ([]*domain.Order, error) {
	return s.repo.AllOrders()
}

func (s *OrderService) UpdateOrderStatus(orderID string, status string) (*domain.Order, string, error) {
	oldStatus := status

	switch status {
	case "0":
		status = "pending"
	case "1":
		status = "shipped"
	case "2":
		status = "delivered"
	case "3":
		status = "canceled"
		// valid status
	default:
		return nil, "", errors.ErrUnsupported
	}
	order, err := s.repo.UpdateOrderStatus(orderID, status)
	if err != nil {
		return nil, "", err
	}
	return order, oldStatus, nil
}
