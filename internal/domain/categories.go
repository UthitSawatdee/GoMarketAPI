package domain


import (
	"time"

	"gorm.io/gorm"
)

// Category represents a product category
type Category struct {
	ID uint `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null;uniqueIndex;size:100"`
	Description string         `json:"description" gorm:"type:text"`
	Products    []Product      `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}


