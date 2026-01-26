package entities

import (
	"time"

	"gorm.io/gorm"
)

// CartItem represents an item in the shopping cart
type CartItem struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CartID    uint           `json:"cart_id" gorm:"index;not null"`
	Cart      Cart           `json:"-" gorm:"foreignKey:CartID"`
	ProductID uint           `json:"product_id" gorm:"index;not null"`
	Product   Product        `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Quantity  int            `json:"quantity" gorm:"not null;default:1"`
	Price     float64        `json:"price" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
