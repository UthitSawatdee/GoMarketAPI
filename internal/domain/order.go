package entities

import (
	"time"

	"gorm.io/gorm"
)

// Order represents a customer order
type Order struct {
	ID uint `json:"id" gorm:"primaryKey"`
	UserID           uint           `json:"user_id" gorm:"index;not null"`
	User             User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Total_amount     float64        `json:"total_amount" gorm:"not null"`
	Status           string         `json:"status" gorm:"size:20;default:pending;index"`
	Shipping_address string         `json:"shipping_address" gorm:"type:text;not null"`
	OrderItems       []OrderItem    `json:"order_items,omitempty" gorm:"foreignKey:OrderID"`
	PaymentMethod    string         `json:"payment_method" gorm:"size:50"`
	PaymentStatus    string         `json:"payment_status" gorm:"size:20;default:pending"`
	Notes            string         `json:"notes" gorm:"type:text"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}
