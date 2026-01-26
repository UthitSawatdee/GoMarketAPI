package container

import (
    "gorm.io/gorm"
	handlers "github.com/Fal2o/E-Commerce_API/internal/adapters/handler"
	adapters "github.com/Fal2o/E-Commerce_API/internal/adapters/repository"
	hash "github.com/Fal2o/E-Commerce_API/pkg/hash"
	usecases "github.com/Fal2o/E-Commerce_API/internal/usecases"
	
)

type Container struct {
    // Handlers
    UserHandler       *handlers.HttpUserHandler
    // ProductHandler    *adapters.HttpProductHandler
    // CategoriesHandler *adapters.HttpCategoriesHandler
    // CartHandler       *adapters.HttpCartHandler
    // OrderHandler      *adapters.HttpOrderHandler
    // HealthHandler     *adapters.HealthHandler
}

func NewContainer(db *gorm.DB) *Container {
    // Repositories
    userRepo := adapters.NewGormUserRepository(db)
    // productRepo := adapters.NewGormProductRepository(db)
    // categoriesRepo := adapters.NewGormCategoriesRepository(db)
    // cartRepo := adapters.NewGormCartRepository(db)
    // orderRepo := adapters.NewGormOrderRepository(db)

    // Services
    passwordService := hash.NewPasswordService()
    userService := usecases.NewUserService(userRepo, passwordService)
    // productService := usecases.NewProductService(productRepo)
    // categoriesService := usecases.NewCategoriesService(categoriesRepo)
    // cartService := usecases.NewCartService(cartRepo)
    // orderService := usecases.NewOrderService(orderRepo, cartRepo, productRepo)

    // Handlers
    return &Container{
        UserHandler:       handlers.NewHttpUserHandler(userService),
        // ProductHandler:    adapters.NewHttpProductHandler(productService),
        // CategoriesHandler: adapters.NewHttpCategoriesHandler(categoriesService),
        // CartHandler:       adapters.NewHttpCartHandler(cartService),
        // OrderHandler:      adapters.NewHttpOrderHandler(orderService),
        // HealthHandler:     adapters.NewHealthHandler(db),
    }
}