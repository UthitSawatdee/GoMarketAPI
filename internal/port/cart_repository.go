package port
import (
	domain "github.com/UthitSawatdee/GoMarketAPI/internal/domain"
)

// CartRepository defines the interface for cart data operations
type CartRepository interface {
	CreateCart(userID *domain.Cart) error
	AddProductToCart(cartItem *domain.CartItem) error
	GetCartByUserID(userID uint) (*domain.Cart, error)
	GetCartItemsByCartID(userID uint) ([]*domain.CartItem, error)
	GetCartItem(cartID uint, productID uint) (*domain.CartItem, error)
	UpdateCartItem(cartItem *domain.CartItem) error
	DeleteProductInCart(cartID uint, productID uint) error
	DeleteAllProductInCart(cartID uint) error
}
