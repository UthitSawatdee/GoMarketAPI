package domain


import (
	"time"

	"gorm.io/gorm"
)

// Cart represents a shopping cart for a user
type Cart struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"uniqueIndex;not null"`
	User      User           `json:"-" gorm:"foreignKey:UserID"`
	Items     []CartItem     `json:"items,omitempty" gorm:"foreignKey:CartID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
