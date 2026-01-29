package handler

import (
	"fmt"
	"strconv"

	usecases "github.com/Fal2o/E-Commerce_API/internal/usecases"
	"github.com/gofiber/fiber/v2"
)

// HttpUserHandler handles HTTP requests for user operations
type HttpCartHandler struct {
	cartUseCase usecases.CartUseCase
}

// NewHttpCartHandler creates a new HttpCartHandler
func NewHttpCartHandler(useCase usecases.CartUseCase) *HttpCartHandler {
	return &HttpCartHandler{cartUseCase: useCase}
}

// AddToCart handles adding items to the cart
// POST /cart/add
func (h *HttpCartHandler) AddProductToCart(c *fiber.Ctx) error {
	productIDStr := c.Params("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}
	fmt.Println("Received product ID to add to cart:", productID)
	if productIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product ID is required",
		})
	}
	userID := c.Locals("user_id").(uint)
	fmt.Println("User ID from context:", userID)
	result, err := h.cartUseCase.AddProductToCart(uint(productID), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add product to cart",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product added to cart successfully",
		"data":    result,
	})

}

func (h *HttpCartHandler) DeleteCartItem(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	productIDStr := c.Params("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}
	fmt.Println("Received product ID to delete from cart:", productID)
	if productIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product ID is required",
		})
	}

	// Call use case to delete cart item
	result, err := h.cartUseCase.DeleteCartItem(uint(productID), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete product from cart",
		})
	}
	if result.Quantity == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Product removed from cart successfully"})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Product deleted from cart successfully",
			"data":    result,
		})
	}
}

func (h *HttpCartHandler) DeleteCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	// Call use case to delete entire cart
	err := h.cartUseCase.DeleteCart(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to clear cart",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Cart cleared successfully",
	})
}

func (h *HttpCartHandler) ViewCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	// Call use case to get cart items
	cartItems, err := h.cartUseCase.ViewCart(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve cart items",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Cart items retrieved successfully",
		"data":    cartItems,
	})
}

func (h *HttpCartHandler) Checkout(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	// Call use case to checkout cart
	order, err := h.cartUseCase.Checkout(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to checkout cart",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Checkout successful",
		"data":    order,
	})	
}