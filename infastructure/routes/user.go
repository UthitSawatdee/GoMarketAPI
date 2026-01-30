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
    user.Get("/cart",c.CartHandler.ViewCart) // get cart 
    user.Post("/cart/item/:product_id",c.CartHandler.AddProductToCart) //add or update product in cart
    user.Delete("/cart/:product_id",c.CartHandler.DeleteCartItem) //decrease or remove product from cart
    user.Delete("/cart/cancel",c.CartHandler.DeleteCart) // cancel cart and all products in cart
    user.Post("/cart/checkout",c.CartHandler.Checkout) // checkout cart (create order and clear cart)

    user.Get("/orders",c.OrderHandler.ViewOrder)
    user.Delete("/order/cancel/:orderID",c.OrderHandler.CancelOrder)




}