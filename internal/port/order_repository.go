package port

import (
	domain "github.com/UthitSawatdee/GoMarketAPI/internal/domain"
)

// OrderRepository defines the interface for order data operations
type OrderRepository interface {
	CreateOrder(userID uint, total_amount float64, orderItem []domain.OrderItem) (uint, error) // Returns created order ID
	GetOrderByUserID(userID uint) ([]*domain.Order, error)
	GetOrderByID(orderID string) (*domain.Order, error) // Get single order with items
	DeleteOrderByOrderID(orderID string) error
	AllOrders() ([]*domain.Order, error)
	UpdateOrderStatus(orderID string, status string) (*domain.Order, error)
}
