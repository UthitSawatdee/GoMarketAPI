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

func (h *HttpCategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")                           // 1. ดึงพารามิเตอร์จาก URL
	request := new(domain.Category)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}
	err := h.CategoryUseCase.UpdateCategory(id,request)
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
