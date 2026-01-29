package routes

import (
    "github.com/Fal2o/E-Commerce_API/infastructure/container"
    "github.com/gofiber/fiber/v2"
)

func setupPublicRoutes(api fiber.Router, c *container.Container) {
    api.Post("/register", c.UserHandler.Register)
    api.Post("/login", c.UserHandler.Login)

}