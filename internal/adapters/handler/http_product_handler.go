package handler

import (
	domain "github.com/UthitSawatdee/GoMarketAPI/internal/domain"
	usecases "github.com/UthitSawatdee/GoMarketAPI/internal/usecases"
	"github.com/gofiber/fiber/v2"
)

type HttpProductHandler struct {
	ProductUseCase usecases.ProductUseCase
}

func NewHttpProductHandler(useCase usecases.ProductUseCase) *HttpProductHandler {
	return &HttpProductHandler{ProductUseCase: useCase}
}

// ProductRequest represents product request body
// @Description Product creation/update request
type ProductRequest struct {
	ID          uint    `json:"id" example:"1"`
	Name        string  `json:"name" example:"iPhone 15 Pro"`
	Description string  `json:"description" example:"Latest Apple smartphone"`
	Price       float64 `json:"price" example:"999.99"`
	Stock       int     `json:"stock" example:"100"`
	CategoryID  uint    `json:"category_id" example:"1"`
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product in the catalog (Admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ProductRequest true "Product details"
// @Success 201 {object} map[string]interface{} "Product created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/product [post]
func (h *HttpProductHandler) CreateProduct(c *fiber.Ctx) error {
	request := new(domain.Product)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	err := h.ProductUseCase.CreateProduct(request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to create product",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Product created successfully",
		"data": fiber.Map{
			"id":          request.ID,
			"name":        request.Name,
			"description": request.Description,
			"price":       request.Price,
			"stock":       request.Stock,
		},
	})
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update an existing product by ID (Admin only)
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Param request body ProductRequest true "Updated product details"
// @Success 200 {object} map[string]interface{} "Product updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/product/{id} [put]
func (h *HttpProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id") // 1. ดึงพารามิเตอร์จาก URL
	request := new(domain.Product)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}
	err := h.ProductUseCase.UpdateProduct(id, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to update product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Product updated successfully",
		"data": fiber.Map{
			"id":          id,
			"name":        request.Name,
			"description": request.Description,
			"price":       request.Price,
			"stock":       request.Stock,
		},
	})
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by ID (Admin only)
// @Tags Products
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]interface{} "Product deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/product/{id} [delete]
func (h *HttpProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.ProductUseCase.DeleteProduct((id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to delete product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Product deleted successfully",
	})
}

// GetAllProducts godoc
// @Summary Get all products
// @Description Retrieve all products from the catalog
// @Tags Products
// @Produce json
// @Success 200 {object} map[string]interface{} "List of products"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /products [get]
func (h *HttpProductHandler) GetAllProducts(c *fiber.Ctx) error {
	data, err := h.ProductUseCase.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to retrieve products",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// GetProductByCategory godoc
// @Summary Get products by category
// @Description Filter products by category name
// @Tags Products
// @Produce json
// @Param category path string true "Category name"
// @Success 200 {object} map[string]interface{} "Filtered products"
// @Failure 404 {object} map[string]interface{} "No products found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /productBy/cat/{category} [get]
func (h *HttpProductHandler) GetProductByCategory(c *fiber.Ctx) error {
	category := c.Params("category")
	data, err := h.ProductUseCase.GetProductByCategory(category)
	if err != nil {
		if err.Error() == "Product not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Product not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to retrieve products by category",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// GetProductByName godoc
// @Summary Search product by name
// @Description Search for a product by its name
// @Tags Products
// @Produce json
// @Param name path string true "Product name"
// @Success 200 {object} map[string]interface{} "Product found"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /product/{name} [get]
func (h *HttpProductHandler) GetProductByName(c *fiber.Ctx) error {
	name := c.Params("name")
	data, err := h.ProductUseCase.GetProductByName(name)
	if err != nil {
		if err.Error() == "Product not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Product not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to retrieve product by name",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}


