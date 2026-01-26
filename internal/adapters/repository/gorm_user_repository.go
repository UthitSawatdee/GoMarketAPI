package repository

import (
	port "github.com/Fal2o/E-Commerce_API/internal/port"
	domain "github.com/Fal2o/E-Commerce_API/internal/domain"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) port.UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(user *domain.User) error {
	if err := r.db.Create(user); err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *GormUserRepository) GetByEmail(email string) (*domain.User, error) {
	user := new(domain.User) // → user เป็น (pointer)
	err := r.db.Where("email =?",email).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}