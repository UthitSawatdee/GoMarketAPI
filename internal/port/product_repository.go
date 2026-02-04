package port

import (
	domain "github.com/UthitSawatdee/GoMarketAPI/internal/domain"
)

// กำหนด "สัญญา" ว่า adapter ต้องทำอะไรได้บ้าง
// คุยกับ gorm
type ProductRepository interface {
	//for adimn
	Create(Product *domain.Product) error
	Update(id string, Product *domain.Product) error
	Delete(id string) error

	//for public
	GetAllProducts() ([]*domain.Product, error)
	GetProductByCategory(category string) ([]*domain.Product, error)
	GetProductByName(Name string) ([]*domain.Product, error) //fiter product by name
	GetProductByID(productID uint) (*domain.Product, error)
	UpdateStock(productID uint, quantity int) error
	RestoreStock(productID uint, quantity int) error // Restore stock when order is cancelled
	// ListUsers() ([]*domain.User, error)
	// GetByEmail(email string) (*domain.User, error)

}
