package usecase

import (
	"errors"
	"fmt"
	domain "github.com/Fal2o/E-Commerce_API/internal/domain"
	port "github.com/Fal2o/E-Commerce_API/internal/port"
	"gorm.io/gorm"
)

// ProductUseCase defines the interface for user business logic
// คุยกับ service (fiber)
type ProductUseCase interface {
	CreateProduct(product *domain.Product) error
	UpdateProduct(id string, product *domain.Product) error
	DeleteProduct(id string) error
	GetAllProducts() ([]*domain.Product, error)
	GetProductByCategory(category string) ([]*domain.Product, error)
	GetProductByName(Name string) ([]*domain.Product, error)
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
	existingProduct, err := s.repo.GetProductByName(product.Name)
	if err != nil {
		return err
	}
	if existingProduct != nil {
		return fmt.Errorf("product name already registered")
	}

	// 3. Create product
	return s.repo.Create(product)
}

func (s *ProductService) UpdateProduct(id string, product *domain.Product) error {
	// Implementation for updating user
	err := s.repo.Update(id, product)
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

func (s *ProductService) GetAllProducts() ([]*domain.Product, error) {
	products, err := s.repo.GetAllProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) GetProductByCategory(category string) ([]*domain.Product, error) {
	products, err := s.repo.GetProductByCategory(category)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) GetProductByName(Name string) ([]*domain.Product, error) {
	product, err := s.repo.GetProductByName(Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil,fmt.Errorf("Product not found")
		}
		return nil, err
	}
	return product, nil
}
