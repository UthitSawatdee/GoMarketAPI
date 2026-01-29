package repository

import (
	domain "github.com/Fal2o/E-Commerce_API/internal/domain"
	port "github.com/Fal2o/E-Commerce_API/internal/port"
	"gorm.io/gorm"
)

type GormOrderRepository struct {
	db *gorm.DB
}

func NewGormOrderRepository(db *gorm.DB) port.OrderRepository {
	return &GormOrderRepository{db: db}
}

func (r *GormOrderRepository) CreateOrder(userID uint, total_amount float64, orderItems []domain.OrderItem) error {
	order := domain.Order{
		UserID:       userID,
		OrderItems:   orderItems,
		Total_amount: total_amount,
	}
	if err := r.db.Create(&order); err.Error != nil {
		return err.Error
	}

	for i := range orderItems {
		orderItems[i].OrderID = order.ID
		orderItems[i].ID = 0
		if err := r.db.Create(&orderItems[i]); err.Error != nil {
			return err.Error
		}
	}

	return nil
}
