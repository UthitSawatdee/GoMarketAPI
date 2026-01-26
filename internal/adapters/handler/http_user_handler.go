package handler

import (
	domain "github.com/Fal2o/E-Commerce_API/internal/domain"
	usecases "github.com/Fal2o/E-Commerce_API/internal/usecases"
	"github.com/gofiber/fiber/v2"
)

type HttpUserHandler struct {
	userUseCase usecases.UserUseCase
}

func NewHttpUserHandler(useCase usecases.UserUseCase) *HttpUserHandler {
	return &HttpUserHandler{userUseCase: useCase}
}

// RegisterRequest represents registration request
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *HttpUserHandler) Register(c *fiber.Ctx) error {
	request := new(RegisterRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "error":   "Invalid request body",
        })
	}

	// Validate request
	if request.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "error":   "Invalid request body",
        })
	}
	if request.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "error":   "Email is required",
        })
	}
	if request.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "error":   "Email is required",
        })
	}
	if len(request.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "error":   "Password must be at least 6 characters",
        })
	}

	user := &domain.User{
		Email:    request.Email,
		Password: request.Password,
		Username: request.Username,
	}

	if err := h.userUseCase.CreateUser(user); err != nil {
		if err.Error() == "email already registered" {
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "success": true,
        "message": "User registered successfully",
        "data": fiber.Map{
            "email":    user.Email,
            "username": user.Username,
        },
    })
}

