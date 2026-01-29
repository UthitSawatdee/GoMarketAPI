package handler

import (
	"test/entities"
	"test/pkg/response"
	"test/usecases"
	"fmt"
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
	request := c.Params("product_id")
	if request == "" {
		return response.BadRequest(c, "product_id is required")
	}
	if request != "" {
		user_id := c.Locals("userID")
		fmt.Println("User ID from context:", user_id)
	} 

	
	err := h.cartUseCase.AddProductToCart(request)
	if err != nil {
		return response.InternalServerError(c, "failed to add product to cart")
	}

	return response.Success(c, "product added to cart successfully", nil)
}