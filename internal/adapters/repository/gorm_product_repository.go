package repository

import (
	"errors"
	domain "github.com/Fal2o/E-Commerce_API/internal/domain"
	port "github.com/Fal2o/E-Commerce_API/internal/port"
	"gorm.io/gorm"
)

type GormProductRepository struct {
	db *gorm.DB
}

func NewGormProductRepository(db *gorm.DB) port.ProductRepository {
	return &GormProductRepository{db: db}
}

func (r *GormProductRepository) Create(product *domain.Product) error {
	if err := r.db.Create(product); err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *GormProductRepository) Update(id string, product *domain.Product) error {
	result := r.db.Where("id = ?", id).Updates(product)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *GormProductRepository) Delete(id string) error {
	result := r.db.Delete(&domain.Product{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *GormProductRepository) GetByName(name string) (*domain.Product, error) {
	product := new(domain.Product) // → product เป็น (pointer)
	err := r.db.Where("name =?", name).First(product).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return product, nil
}
