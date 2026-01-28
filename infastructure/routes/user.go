package routes

import (
    "github.com/Fal2o/E-Commerce_API/infastructure/container"
    "github.com/gofiber/fiber/v2"
    "github.com/Fal2o/E-Commerce_API/infastructure/config"
    "github.com/Fal2o/E-Commerce_API/internal/middleware"


)

func setupUserRoutes(api fiber.Router, c *container.Container,cfg *config.Config) {
        admin := api.Group("/user",
        middleware.AuthMiddleware(cfg.JWT.Secret),
        middleware.UserOnly(),
    )
	// Update user profile
    admin.Put("/profile", c.UserHandler.UpdateProfile)
	// Get user profile
    admin.Get("/profile", c.UserHandler.GetProfile)


}