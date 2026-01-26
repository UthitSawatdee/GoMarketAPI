package domain


import (
	"time"

	"gorm.io/gorm"
)

// OrderItem represents an item within an order
type OrderItem struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	OrderID     uint           `json:"order_id" gorm:"index;not null"`
	Order       Order          `json:"-" gorm:"foreignKey:OrderID"`
	ProductID   uint           `json:"product_id" gorm:"index;not null"`
	Product     Product        `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	ProductName string         `json:"product_name" gorm:"size:255"`
	Quantity    int            `json:"quantity" gorm:"not null"`
	Price       float64        `json:"price" gorm:"not null"`
	Subtotal    float64        `json:"subtotal" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
