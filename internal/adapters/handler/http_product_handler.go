package handler

import (
	domain "github.com/Fal2o/E-Commerce_API/internal/domain"
	usecases "github.com/Fal2o/E-Commerce_API/internal/usecases"
	"github.com/gofiber/fiber/v2"
)

type HttpProductHandler struct {
	ProductUseCase usecases.ProductUseCase
}

func NewHttpProductHandler(useCase usecases.ProductUseCase) *HttpProductHandler {
	return &HttpProductHandler{ProductUseCase: useCase}
}

type ProductRequest struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

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
