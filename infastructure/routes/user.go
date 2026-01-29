package routes

import (
    "github.com/Fal2o/E-Commerce_API/infastructure/container"
    "github.com/gofiber/fiber/v2"
    "github.com/Fal2o/E-Commerce_API/infastructure/config"
    "github.com/Fal2o/E-Commerce_API/internal/middleware"


)

func setupUserRoutes(api fiber.Router, c *container.Container,cfg *config.Config) {
        user := api.Group("/user",
        middleware.AuthMiddleware(cfg.JWT.Secret),
        middleware.UserOnly(),
    )
	// Update user profile
    user.Put("/profile", c.UserHandler.UpdateProfile)
	// Get user profile
    user.Get("/profile", c.UserHandler.GetProfile)

    // Cart routes
    // user.Get("/cart",c.UserHandler.ViewCart)
    user.Put("/cart/:product_id",c.CartHandler.AddProductToCart)
    // user.Delete("/cart/:id",c.UserHandler.RemoveProduct)
    // user.Delete("/cart",c.UserHandler.DeleteCart)
    // user.Post("/cart/checkout",c.UserHandler.Checkout)

    // user.Get("/orders",c.UserHandler.ViewOrder)
    // user.Get("/orders/:id",c.UserHandler.ViewOrderDetail)
    // user.Patch("/order/:id/cancel",c.UserHandler.CancelOrder)




}