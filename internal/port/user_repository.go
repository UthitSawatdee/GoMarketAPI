package port

import (
	domain "github.com/UthitSawatdee/GoMarketAPI/internal/domain"
)

// กำหนด "สัญญา" ว่า repository ต้องทำอะไรได้บ้าง
type UserRepository interface {
	Create(user *domain.User) error
	// Update(user *domain.User) error
	// Delete(id uint) error
	GetUserByEmail(email string) (*domain.User, error)
	GetUserByID(id uint) (*domain.User, error)	
	// ListUsers() ([]*domain.User, error)
	// GetByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	AllUsers() ([]*domain.User, error)
}
