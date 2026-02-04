package usecase

import (
	"errors"
	"fmt"

	"github.com/UthitSawatdee/GoMarketAPI/internal/domain"
	port "github.com/UthitSawatdee/GoMarketAPI/internal/port"
)

// OrderUseCase defines the interface for order business logic
type OrderUseCase interface {
	ViewOrder(userID uint) ([]*domain.Order, error)
	CancelOrder(orderID string, userID uint) error
	AllOrders() ([]*domain.Order, error)
	UpdateOrderStatus(orderID string, status string) (*domain.Order, string, error)
}

type OrderService struct {
	repo        port.OrderRepository
	productRepo port.ProductRepository
}

// NewOrderService creates a new OrderService with required dependencies
func NewOrderService(repo port.OrderRepository, productRepo port.ProductRepository) OrderUseCase {
	return &OrderService{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (s *OrderService) ViewOrder(userID uint) ([]*domain.Order, error) {
	order, err := s.repo.GetOrderByUserID(userID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// CancelOrder cancels an order and restores product stock
func (s *OrderService) CancelOrder(orderID string, userID uint) error {
	// 1. Get order with items to restore stock
	order, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// 2. Verify ownership (user can only cancel their own orders)
	if order.UserID != userID {
		return errors.New("unauthorized: cannot cancel another user's order")
	}

	// 3. Check if order can be cancelled
	currentStatus := domain.OrderStatus(order.Status)
	if !currentStatus.IsCancellable() {
		return fmt.Errorf("order cannot be cancelled: current status is %s", order.Status)
	}

	// 4. Restore stock for each item in the order
	for _, item := range order.OrderItems {
		if err := s.productRepo.RestoreStock(item.ProductID, item.Quantity); err != nil {
			// Log the error but continue with other items
			// In production, this should be wrapped in a transaction
			fmt.Printf("Warning: failed to restore stock for product %d: %v\n", item.ProductID, err)
		}
	}

	// 5. Update order status to cancelled instead of deleting
	_, err = s.repo.UpdateOrderStatus(orderID, string(domain.OrderStatusCancelled))
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}

func (s *OrderService) AllOrders() ([]*domain.Order, error) {
	return s.repo.AllOrders()
}

// UpdateOrderStatus validates and updates order status with proper transitions
func (s *OrderService) UpdateOrderStatus(orderID string, newStatus string) (*domain.Order, string, error) {
	// 1. Get current order to check current status
	order, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		return nil, "", fmt.Errorf("order not found: %w", err)
	}

	// 2. Parse and validate status
	targetStatus := domain.OrderStatus(newStatus)
	if !targetStatus.IsValid() {
		// Support legacy numeric status codes for backward compatibility
		switch newStatus {
		case "0":
			targetStatus = domain.OrderStatusPending
		case "1":
			targetStatus = domain.OrderStatusShipped
		case "2":
			targetStatus = domain.OrderStatusDelivered
		case "3":
			targetStatus = domain.OrderStatusCancelled
		default:
			return nil, "", fmt.Errorf("invalid status: %s", newStatus)
		}
	}

	// 3. Validate transition
	currentStatus := domain.OrderStatus(order.Status)
	if !currentStatus.CanTransitionTo(targetStatus) {
		return nil, "", fmt.Errorf("invalid transition: cannot change from '%s' to '%s'", currentStatus, targetStatus)
	}

	// 4. If cancelling, restore stock
	if targetStatus == domain.OrderStatusCancelled {
		for _, item := range order.OrderItems {
			if err := s.productRepo.RestoreStock(item.ProductID, item.Quantity); err != nil {
				fmt.Printf("Warning: failed to restore stock for product %d: %v\n", item.ProductID, err)
			}
		}
	}

	// 5. Update status
	updatedOrder, err := s.repo.UpdateOrderStatus(orderID, string(targetStatus))
	if err != nil {
		return nil, "", err
	}

	return updatedOrder, string(currentStatus), nil
}
