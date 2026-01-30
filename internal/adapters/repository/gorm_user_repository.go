package repository

import (
	"errors"
	domain "github.com/Fal2o/E-Commerce_API/internal/domain"
	port "github.com/Fal2o/E-Commerce_API/internal/port"
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

func (r *GormUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	user := new(domain.User) // → user เป็น (pointer)
	err := r.db.Where("email =?", email).First(user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *GormUserRepository) GetUserByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *GormUserRepository) Update(user *domain.User) error {
	result := r.db.Where("id = ?", user.ID).Updates(user) 
	if result.Error != nil {
		return result.Error
	}
    if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }
	return nil
}

func (r *GormUserRepository) AllUsers() ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
