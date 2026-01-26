package port

import (
	domain "github.com/Fal2o/E-Commerce_API/internal/domain"
)

// กำหนด "สัญญา" ว่า repository ต้องทำอะไรได้บ้าง
type UserRepository interface {
	Create(user *domain.User) error
	// Update(user *domain.User) error
	// Delete(id uint) error
	GetByEmail(email string) (*domain.User, error)
	// GetUser(id uint) (*domain.User, error)
	// ListUsers() ([]*domain.User, error)
	// GetByEmail(email string) (*domain.User, error)
}
