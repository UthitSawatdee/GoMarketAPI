package routes

import (
    "github.com/Fal2o/E-Commerce_API/infrastructure/container"
    "github.com/gofiber/fiber/v2"
    "github.com/Fal2o/E-Commerce_API/infrastructure/config"
    "github.com/Fal2o/E-Commerce_API/internal/middleware"


)

func setupAdminRoutes(api fiber.Router, c *container.Container,cfg *config.Config) {
        admin := api.Group("/admin",
        middleware.AuthMiddleware(cfg.JWT.Secret),
        middleware.AdminOnly(),
    )
	
	admin.Post("/product", c.ProductHandler.CreateProduct)
    admin.Put("/product/:id", c.ProductHandler.UpdateProduct)
    admin.Delete("/product/:id", c.ProductHandler.DeleteProduct)

    admin.Post("/category", c.CategoriesHandler.CreateCategory)
    admin.Put("/category/:id", c.CategoriesHandler.UpdateCategory)
    admin.Delete("/category/:id", c.CategoriesHandler.DeleteCategory)

    admin.Get("/users", c.UserHandler.AllUsers)
    admin.Get("/orders", c.OrderHandler.ViewAllOrders)
    admin.Put("/order/status/:orderID/:status", c.OrderHandler.UpdateOrderStatus)
}