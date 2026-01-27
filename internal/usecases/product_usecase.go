package usecase

import (
	"fmt"
	domain "github.com/Fal2o/E-Commerce_API/internal/domain"
	port "github.com/Fal2o/E-Commerce_API/internal/port"
	"gorm.io/gorm"
	"errors"

)

// ProductUseCase defines the interface for user business logic
type ProductUseCase interface {
	CreateProduct(product *domain.Product) error
	UpdateProduct(id string,product *domain.Product) error
	DeleteProduct(id string) error
}

type ProductService struct {
	repo port.ProductRepository
}

func NewProductService(repo port.ProductRepository) ProductUseCase {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) CreateProduct(product *domain.Product) error {
	// 1. Check if email already exists
	existingProduct, err := s.repo.GetByName(product.Name)
	if err != nil {
		return err
	}
	if existingProduct != nil {
		return fmt.Errorf("product name already registered")
	}

	// 3. Create product
	return s.repo.Create(product)
}

func (s *ProductService) UpdateProduct(id string,product *domain.Product) error {
	// Implementation for updating user
	err := s.repo.Update(id,product)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("Product not found")
		}
		return err
	}
	return nil 
}

func (s *ProductService) DeleteProduct(id string) error {
	// Implementation for deleting user
	err := s.repo.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("Product not found")
		}
		return err
	}
	return nil
}