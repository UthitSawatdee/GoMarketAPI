package port

import (
	domain "github.com/UthitSawatdee/GoMarketAPI/internal/domain"
)

// OrderRepository defines the interface for order data operations
type OrderRepository interface {
	CreateOrder(userID uint, total_amount float64,orderItem []domain.OrderItem) error
	GetOrderByUserID(userID uint) ([]*domain.Order, error)
	DeleteOrderByOrderID(orderID string) error
	AllOrders() ([]*domain.Order, error)
	UpdateOrderStatus(orderID string, status string) (*domain.Order, error)
}
