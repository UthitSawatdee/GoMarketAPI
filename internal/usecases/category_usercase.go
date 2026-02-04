package usecase

import (
	"fmt"
	domain "github.com/UthitSawatdee/GoMarketAPI/internal/domain"
	port "github.com/UthitSawatdee/GoMarketAPI/internal/port"
	"gorm.io/gorm"
	"errors"
)

// ProductUseCase defines the interface for user business logic
type CategoryUseCase interface {
	CreateCategory(category *domain.Category) error
	UpdateCategory(id string, category *domain.Category) error
	DeleteCategory(id string) error
}

type CategoryService struct {
	repo port.CategoryRepository
}

func NewCategoryService(repo port.CategoryRepository) CategoryUseCase {
	return &CategoryService{
		repo: repo,
	}
}

func (s *CategoryService) CreateCategory(category *domain.Category) error {
	// 1. Check if email already exists
	existingCategory, err := s.repo.GetByName(category.Name)
	if err != nil {
		return err
	}
	if existingCategory != nil {
		return fmt.Errorf("Category name already exited")
	}

	// 3. Create product
	return s.repo.Create(category)
}

func (s *CategoryService) UpdateCategory(id string, category *domain.Category) error {
	// Implementation for updating user
	if err := s.repo.Update(id, category); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("Category not found")
		}
		return err
	}
	return nil
}

func (s *CategoryService) DeleteCategory(id string) error {
	// Implementation for deleting user
	existingCategory, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existingCategory == nil {
		return fmt.Errorf("Category not found")
	}
	return s.repo.Delete(id)
}
