package repository

import (
	domain "github.com/Fal2o/E-Commerce_API/internal/domain"
	"gorm.io/gorm"
)

type GormCartRepository struct {
	db *gorm.DB
}

func NewGormCartRepository(db *gorm.DB) *GormCartRepository {
	return &GormCartRepository{db:db}
}

// Implement CartRepository methods here
func (r *GormCartRepository) GetCartByUserID(userID uint) (*domain.Cart, error) {
	cart := new(domain.Cart)
	err := r.db.Where("user_id =?", userID).First(cart).Error
	if err != nil {
		return nil,err
	}
	return cart,nil
}

func (r *GormCartRepository) GetCartItem(cartID uint, productID uint) (*domain.CartItem, error) {
	cartItem := new(domain.CartItem)
	err := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).First(cartItem).Error
	if err != nil {
		return nil, err
	}
	return cartItem, nil
}

func (r *GormCartRepository) AddProductToCart(cartItem *domain.CartItem) error {
	if err := r.db.Create(cartItem); err.Error != nil {
		return err.Error
	}
	return nil
}
 

func (r *GormCartRepository) CreateCart(userID *domain.Cart) error {
	if err := r.db.Create(userID); err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *GormCartRepository) UpdateCartItem(cartItem *domain.CartItem) error {
	result := r.db.Model(&domain.CartItem{}).Where("id = ?", cartItem.ID).Updates(cartItem)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}	
	return nil
}

func (r *GormCartRepository) DeleteProductInCart(cartID uint,productID uint) error {
	result := r.db.Where("cart_id = ? AND product_id = ?", cartID,productID).Delete(&domain.CartItem{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormCartRepository) DeleteAllProductInCart(cartID uint) error {
	result := r.db.Where("cart_id = ?", cartID).Delete(&domain.CartItem{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormCartRepository) GetCartItemsByCartID(userID uint) ([]*domain.CartItem, error) {
	var cartItems []*domain.CartItem
	err := r.db.Where("cart_id = ?", userID).Find(&cartItems).Error
	if err != nil {
		return nil, err
	}
	return cartItems, nil
}