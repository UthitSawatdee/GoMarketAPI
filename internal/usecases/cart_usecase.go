package usecase

import (
	"errors"

	"github.com/UthitSawatdee/GoMarketAPI/internal/domain"
	port "github.com/UthitSawatdee/GoMarketAPI/internal/port"
	"gorm.io/gorm"
)

// CartUseCase defines the interface for cart business logic
type CartUseCase interface {
	AddProductToCart(productID uint, userID uint) (*CartItemResult, error)
	DeleteCartItem(cartID uint, productID uint) (*CartItemResult, error)
	DeleteCart(userID uint) error
	ViewCart(userID uint) ([]*CartItemResult, error)
	Checkout(userID uint) ([]domain.OrderItem, error)
}

type CartService struct {
	repo        port.CartRepository
	productRepo port.ProductRepository
	orderRepo   port.OrderRepository
}

func NewCartService(repo port.CartRepository, productRepo port.ProductRepository, orderRepo port.OrderRepository) CartUseCase {
	return &CartService{
		repo:        repo,
		productRepo: productRepo,
		orderRepo:   orderRepo,
	}
}

// CartItemResult represents the result of cart operations
type CartItemResult struct {
	ProductName string
	Quantity    int
	UnitPrice   float64
	TotalPrice  float64
}

func (s *CartService) AddProductToCart(productID uint, userID uint) (*CartItemResult, error) {
	// Business logic to add product to cart
	// 1. ดึงข้อมูล Product เพื่อเอา Price
	product, err := s.productRepo.GetProductByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Check stock availability
	if !product.HasStock(product.Stock) {
		return nil, errors.New("product out of stock")
	}

	// 3.1: Get or Create Cart
	cart, err := s.getOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	// 3.2: Add or Update CartItem
	cartItem, err := s.addOrUpdateCartItem(cart.ID, productID, 1, product.Price)
	if err != nil {
		return nil, err
	}

	// 3.3: Update Product Stock (ลด stock)
	if err := s.productRepo.UpdateStock(productID, 1); err != nil {
		return nil, err
	}

	// 3.4: Build result
	result := &CartItemResult{
		ProductName: product.Name,
		Quantity:    cartItem.Quantity,
		UnitPrice:   product.Price,
		TotalPrice:  float64(cartItem.Quantity) * product.Price,
	}
	return result, nil
}

// getOrCreateCart retrieves existing cart or creates a new one
func (s *CartService) getOrCreateCart(userID uint) (*domain.Cart, error) {
	cart, err := s.repo.GetCartByUserID(userID)
	if err != nil {
		// ถ้าไม่มี Cart → สร้างใหม่
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newCart := &domain.Cart{UserID: userID}
			if err := s.repo.CreateCart(newCart); err != nil {
				return nil, err
			}
			return newCart, nil
		}
		return nil, err
	}
	return cart, nil
}

// addOrUpdateCartItem adds new item or updates existing item quantity
func (s *CartService) addOrUpdateCartItem(
	cartID uint,
	productID uint,
	quantity int,
	price float64,
) (*domain.CartItem, error) {
	// Try to get existing item
	existingItem, err := s.repo.GetCartItem(cartID, productID)
	if err != nil {
		// ถ้าไม่มี Item → เพิ่มใหม่
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newItem := &domain.CartItem{
				CartID:    cartID,
				ProductID: productID,
				Quantity:  quantity,
				Price:     price,
			}
			if err := s.repo.AddProductToCart(newItem); err != nil {
				return nil, err
			}
			return newItem, nil
		}
		return nil, err
	}

	// มี Item อยู่แล้ว → Update quantity
	existingItem.Quantity += quantity
	existingItem.Price *= float64(existingItem.Quantity) // อัพเดทราคาล่าสุด

	if err := s.repo.UpdateCartItem(existingItem); err != nil {
		return nil, err
	}

	return existingItem, nil
}

func (s *CartService) DeleteCartItem(productID uint, userID uint) (*CartItemResult, error) {
	product, err := s.productRepo.GetProductByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Check stock availability
	if !product.HasStock(product.Stock) {
		return nil, errors.New("product out of stock")
	}

	// 3.1: Get or Create Cart
	cart, err := s.getOrCreateCart(userID)
	if err != nil {
		return nil, err
	}

	cartItem, err := s.repo.GetCartItem(cart.ID, productID)
	if err != nil {
		return nil, err
	}

	// 3.2: Remove CartItem
	cartItem.Quantity -= 1
	err = s.repo.UpdateCartItem(cartItem)
	if err != nil {
		return nil, err
	}
	if cartItem.Quantity != 0 {
		// 3.3: Update Product Stock (เพิ่ม stock)
		if err := s.productRepo.UpdateStock(productID, -1); err != nil {
			return nil, err
		}
	} else {
		// ถ้า quantity = 0 → ลบ item ออกจาก cart
		err = s.repo.DeleteProductInCart(cartItem.CartID, cartItem.ProductID)
		if err != nil {
			return nil, err
		}
	}

	// 3.4: Build result
	result := &CartItemResult{
		ProductName: product.Name,
		Quantity:    cartItem.Quantity,
		UnitPrice:   product.Price,
		TotalPrice:  float64(cartItem.Quantity) * product.Price,
	}
	return result, nil
}

func (s *CartService) DeleteCart(userID uint) error {
	// 1: Get Cart
	cart, err := s.repo.GetCartByUserID(userID)
	if err != nil {
		return err
	}

	// 2: Clear Cart
	err = s.repo.DeleteAllProductInCart(cart.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *CartService) ViewCart(userID uint) ([]*CartItemResult, error) {
	// 1: Get Cart
	cart, err := s.repo.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 2: Get Cart Items
	cartItems, err := s.repo.GetCartItemsByCartID(cart.ID)
	if err != nil {
		return nil, err
	}

	// 3: Build Results
	var results []*CartItemResult
	for _, item := range cartItems {
		product, err := s.productRepo.GetProductByID(item.ProductID)
		if err != nil {
			return nil, err
		}
		result := &CartItemResult{
			ProductName: product.Name,
			Quantity:    item.Quantity,
			UnitPrice:   product.Price,
			TotalPrice:  float64(item.Quantity) * product.Price,
		}
		results = append(results, result)
	}
	return results, nil
}

func (s *CartService) Checkout(userID uint) ([]domain.OrderItem, error) {
	// 1: Get Cart
	cart, err := s.repo.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 2: Get Cart Items
	cartItems, err := s.repo.GetCartItemsByCartID(cart.ID)
	if err != nil {
		return nil, err
	}

	// 3: Build Results
	var results []domain.OrderItem
	var totalAmount float64
	for _, item := range cartItems {
		product, err := s.productRepo.GetProductByID(item.ProductID)
		if err != nil {
			return nil, err
		}
		itemSubtotal := float64(item.Quantity) * product.Price
		result := &domain.OrderItem{
			ProductID:   product.ID, // bug ก่อนหน้านี้ไม่ใ่ส
			ProductName: product.Name,
			Quantity:    item.Quantity,
			Price:       product.Price,
			Subtotal:    itemSubtotal,
		}
		totalAmount += itemSubtotal
		results = append(results, *result)
	}
	if err_order := s.orderRepo.CreateOrder(userID, totalAmount, results); err_order != nil {
		return nil, err_order
	}

	// 4: Clear Cart after checkout
	err = s.repo.DeleteAllProductInCart(cart.ID)
	if err != nil {
		return nil, err
	}

	return results, nil
}
