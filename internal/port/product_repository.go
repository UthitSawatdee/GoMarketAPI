package port

import (
	domain "github.com/Fal2o/E-Commerce_API/internal/domain"
)

// กำหนด "สัญญา" ว่า repository ต้องทำอะไรได้บ้าง
type ProductRepository interface {
	Create(Product *domain.Product) error
	Update(id string,Product *domain.Product) error
	Delete(id string) error
	GetByName(Name string) (*domain.Product, error)
	// GetUser(id uint) (*domain.User, error)
	// ListUsers() ([]*domain.User, error)
	// GetByEmail(email string) (*domain.User, error)
	
}
