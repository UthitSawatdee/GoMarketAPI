package routes

import (
	"github.com/UthitSawatdee/GoMarketAPI/infrastructure/container"
	"github.com/gofiber/fiber/v2"
)

func setupPublicRoutes(api fiber.Router, c *container.Container) {
	api.Post("/register", c.UserHandler.Register)
	api.Post("/login", c.UserHandler.Login)
	//Search & Filter by Category
	api.Get("/products", c.ProductHandler.GetAllProducts)
	api.Get("/product/:name", c.ProductHandler.GetProductByName)
	api.Get("/productBy/cat/:category", c.ProductHandler.GetProductByCategory)

}
