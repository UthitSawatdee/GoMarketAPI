package handler

import (
	domain "github.com/Fal2o/E-Commerce_API/internal/domain"
	usecases "github.com/Fal2o/E-Commerce_API/internal/usecases"
	"github.com/gofiber/fiber/v2"
)

type HttpCategoryHandler struct {
	CategoryUseCase usecases.CategoryUseCase
}

func NewHttpCategoryHandler(useCase usecases.CategoryUseCase) *HttpCategoryHandler {
	return &HttpCategoryHandler{CategoryUseCase: useCase}
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new product category (Admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.Category true "Category details"
// @Success 201 {object} map[string]interface{} "Category created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/category [post]
func (h *HttpCategoryHandler) CreateCategory(c *fiber.Ctx) error {
	request := new(domain.Category)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	err := h.CategoryUseCase.CreateCategory(request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to create product",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Category created successfully",
		"data": fiber.Map{
			"id":   request.ID,
			"name": request.Name,
		},
	})
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update an existing category by ID (Admin only)
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Param request body domain.Category true "Updated category details"
// @Success 200 {object} map[string]interface{} "Category updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/category/{id} [put]
func (h *HttpCategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id") // 1. ดึงพารามิเตอร์จาก URL
	request := new(domain.Category)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}
	err := h.CategoryUseCase.UpdateCategory(id, request)
	if err != nil {
		if err.Error() == "Category not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Category not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to update category",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Category updated successfully",
		"data": fiber.Map{
			"id":   id,
			"name": request.Name,
		},
	})
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a category by ID (Admin only)
// @Tags Categories
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} map[string]interface{} "Category deleted successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Admin access required"
// @Failure 404 {object} map[string]interface{} "Category not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/category/{id} [delete]
func (h *HttpCategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.CategoryUseCase.DeleteCategory(id)
	if err != nil {
		if err.Error() == "Category not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"error":   "Category not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to delete category",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Category deleted successfully",
	})
}
