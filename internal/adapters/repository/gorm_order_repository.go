package repository

import (
	domain "github.com/UthitSawatdee/GoMarketAPI/internal/domain"
	port "github.com/UthitSawatdee/GoMarketAPI/internal/port"
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
        OrderItems:   orderItems,  // GORM จะสร้าง orderItems ให้เอง
        Total_amount: total_amount,
    }
    
    // สร้างทั้ง order และ orderItems ในครั้งเดียว
    if err := r.db.Create(&order).Error; err != nil {
        return err
    }
    
    return nil
}

func (r *GormOrderRepository) GetOrderByUserID(userID uint) ([]*domain.Order, error) {
	var order []*domain.Order
	if err := r.db.Preload("OrderItems").Where("user_id = ?", userID).Find(&order); err.Error != nil {
		return nil, err.Error
	}
	return order, nil
}

func (r *GormOrderRepository) DeleteOrderByOrderID(orderID string) error {
	resultOrder := r.db.Where("id = ?", orderID).Delete(&domain.Order{})
	if resultOrder.Error != nil {
		return resultOrder.Error
	}

	resultOrderItems := r.db.Where("order_id = ?", orderID).Delete(&domain.OrderItem{})
	if resultOrderItems.Error != nil {
		return resultOrderItems.Error
	}

	if resultOrder.RowsAffected == 0 || resultOrderItems.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	
	return nil
}

func (r *GormOrderRepository) AllOrders() ([]*domain.Order, error) {
	var orders []*domain.Order
	err := r.db.Preload("OrderItems").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *GormOrderRepository) UpdateOrderStatus(orderID string, status string) (*domain.Order, error) {
	var order domain.Order
	result := r.db.Model(&order).Where("id = ?", orderID).Update("status", status)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &order, nil
}	