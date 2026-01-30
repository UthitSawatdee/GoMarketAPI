package routes

import (
    "github.com/Fal2o/E-Commerce_API/infrastructure/container"
    "github.com/Fal2o/E-Commerce_API/infrastructure/config"
    "github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, c *container.Container, cfg *config.Config) {
    // API v1
    api := app.Group("/api/v1")

    setupPublicRoutes(api, c)
    setupUserRoutes(api, c, cfg)
    setupAdminRoutes(api, c, cfg)
}
