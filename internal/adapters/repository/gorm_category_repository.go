package repository

import (
	"errors"
	domain "github.com/Fal2o/E-Commerce_API/internal/domain"
	port "github.com/Fal2o/E-Commerce_API/internal/port"
	"gorm.io/gorm"
)

type GormCategoryRepository struct {
	db *gorm.DB
}

func NewGormCategoryRepository(db *gorm.DB) port.CategoryRepository {
	return &GormCategoryRepository{db: db}
}

func (r *GormCategoryRepository) Create(category *domain.Category) error {
	if err := r.db.Create(category); err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *GormCategoryRepository) Update(id string, category *domain.Category) error {
	result := r.db.Where("id = ?", id).Updates(category) 
	if result.Error != nil {
		return result.Error
	}
    if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }
	return nil
}

func (r *GormCategoryRepository) Delete(id string) error {
	result := r.db.Delete(&domain.Category{}, id)
	if result.Error != nil {
		return result.Error
	}
	// 3. เช็คว่าลบได้กี่แถว
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *GormCategoryRepository) GetByName(name string) (*domain.Category, error) {
	category := new(domain.Category)
	err := r.db.Where("name =?", name).First(category).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *GormCategoryRepository) GetByID(id string) (*domain.Category, error) {
	category := new(domain.Category)
	err := r.db.Where("id =?", id).First(category).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return category, nil
}
