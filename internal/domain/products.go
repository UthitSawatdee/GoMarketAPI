package domain

import (
	"time"
	"errors"
	"gorm.io/gorm"
)

// Product represents a product in the catalog
type Product struct {
	ID          uint     `json:"id" gorm:"primaryKey"`
	Name        string   `json:"name" gorm:"not null;size:255;index"`
	Description string   `json:"description" gorm:"type:text"`
	Price       float64  `json:"price" gorm:"not null;default:0"`
	Stock       int      `json:"stock" gorm:"not null;default:0"`
	CategoryID  uint     `json:"category_id" gorm:"index;default:0"`
	Category    Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	// IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// HasStock checks if product has enough stock
func (p *Product) HasStock(quantity int) bool {
    return p.Stock >= quantity
}

// DeductStock reduces the stock by given quantity
func (p *Product) DeductStock(quantity int) error {
    if !p.HasStock(quantity) {
        return errors.New("insufficient stock")
    }
    p.Stock -= quantity
    return nil
}