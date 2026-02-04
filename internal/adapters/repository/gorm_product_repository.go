package repository

import (
	"errors"
	"fmt"

	domain "github.com/UthitSawatdee/GoMarketAPI/internal/domain"
	port "github.com/UthitSawatdee/GoMarketAPI/internal/port"
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

func (r *GormProductRepository) GetAllProducts() ([]*domain.Product, error) {
	var products []*domain.Product
    err := r.db.Preload("Category").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *GormProductRepository) GetProductByCategory(category string) ([]*domain.Product, error) {
	var products []*domain.Product
	err := r.db.Preload("Category").Where("category_id = ? ", category).Find(&products)
	if err.Error != nil {
		return nil, err.Error
	}
	if err.RowsAffected == 0 {
		return nil, fmt.Errorf("Product not found")
	}

	return products, nil
}

func (r *GormProductRepository) GetProductByName(name string) ([]*domain.Product, error) {
	var product []*domain.Product
	err := r.db.Preload("Category").Where("LOWER(name) LIKE LOWER(?)", "%"+name+"%").Find(&product)
	if err.Error != nil {
		return nil, err.Error
	}

	if err.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return product, nil
}

func (r *GormProductRepository) GetProductByID(productID uint) (*domain.Product, error) {
	product := new(domain.Product)
	err := r.db.Preload("Category").First(product, productID).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

// UpdateStock updates product stock with optimistic locking
// quantity: จำนวนที่จะลด (ค่าบวก = ลด, ค่าลบ = เพิ่ม)
func (r *GormProductRepository) UpdateStock(productID uint, quantity int) error {
    // ใช้ raw SQL เพื่อ atomic update และป้องกัน race condition
    result := r.db.Model(&domain.Product{}).
        Where("id = ? AND stock >= ?", productID, quantity).
        Update("stock", gorm.Expr("stock - ?", quantity))

    if result.Error != nil {
        return fmt.Errorf("%w: %v", errors.New("database error"), result.Error)
    }

    // ถ้า RowsAffected = 0 หมายความว่า stock ไม่พอ หรือ product ไม่มี
    if result.RowsAffected == 0 {
        return errors.New("insufficient stock")
    }

    return nil
}
