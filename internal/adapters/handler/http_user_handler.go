package handler

import (
	"fmt"

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
// @Description User registration request body
type RegisterRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"securepass123"`
	Username string `json:"username" example:"johndoe"`
	Role     string `json:"role" example:"user" enums:"user,admin"`
}

// LoginRequest represents login request
// @Description User login request body
type LoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"securepass123"`
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with email, password, and username
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration details"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body or validation error"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /register [post]
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
		Role:     request.Role,
	}

	if err := h.userUseCase.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to register user",
			"error":   err.Error(),
		})
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

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password, returns JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "User login credentials"
// @Success 200 {object} map[string]interface{} "Login successful with JWT token"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Router /login [post]
func (h *HttpUserHandler) Login(c *fiber.Ctx) error {
	request := new(LoginRequest)
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
			"error":   "Email is required",
		})
	}
	if request.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Password is required",
		})
	}

	user := &domain.User{
		Email:    request.Email,
		Password: request.Password,
	}

	// Call use case to handle login
	err, token := h.userUseCase.LoginUser(user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Login failed",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Login successful",
		"data": fiber.Map{
			"token": token,
		},
	})
}

func (h *HttpUserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	// fmt.Println("Retrieved userID from context:", userID)
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "userID not found",
		})
	}
	userIDUint, ok := userID.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "invalid userID",
		})
	}
	user, err := h.userUseCase.GetUserByID(userIDUint)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to retrieve user profile",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "User profile retrieved successfully",
		"data": fiber.Map{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// UpdateProfileRequest represents update profile request
type UpdateProfileRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}

// UpdateProfile updates the current user's profile
// PUT /user/profile
func (h *HttpUserHandler) UpdateProfile(c *fiber.Ctx) error {
	userIDValue := c.Locals("user_id")
	if userIDValue == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "userID not found",
		})
	}
	userID, ok := userIDValue.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "invalid userID",
		})
	}

	request := new(UpdateProfileRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}
	// Get existing user
	user, err := h.userUseCase.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
			"error":   err.Error(),
		})
	}
	if request.Username != "" {
		user.Username = request.Username
	}
	if request.Email != "" {
		user.Email = request.Email
	}
	fmt.Println("Request Password:", request.Password)
	fmt.Println("Request NewPassword:", request.NewPassword)

	if err := h.userUseCase.UpdateUser(user, request.Password, request.NewPassword); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to update user profile",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":  true,
		"ID":       user.ID,
		"Email":    user.Email,
		"Username": user.Username,
	})
}

func (h *HttpUserHandler) AllUsers(c *fiber.Ctx) error {
	users, err := h.userUseCase.AllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to retrieve users",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Users retrieved successfully",
		"data":    users,
	})
}
