package port

import (
	domain "github.com/Fal2o/E-Commerce_API/internal/domain"
)

// กำหนด "สัญญา" ว่า repository ต้องทำอะไรได้บ้าง
type CategoryRepository interface {
	Create(Category *domain.Category) error
	Update(id string,Category *domain.Category) error
	Delete(id string) error
	GetByName(Name string) (*domain.Category, error)
	GetByID(id string) (*domain.Category, error)
	// GetUser(id uint) (*domain.User, error)
	// ListUsers() ([]*domain.User, error)
	// GetByEmail(email string) (*domain.User, error)
	
}
