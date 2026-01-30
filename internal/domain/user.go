package domain

import (
	"time"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null;size:255"`
	Password  string         `json:"-" gorm:"not null"` // password จะไม่ถูกส่งกลับใน JSON
	Username  string         `json:"username" gorm:"not null;size:100"`
	Role      string         `json:"role" gorm:"default:customer;size:20"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
