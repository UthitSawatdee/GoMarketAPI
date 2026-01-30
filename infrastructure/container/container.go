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
    ProductHandler    *handlers.HttpProductHandler
    CategoriesHandler *handlers.HttpCategoryHandler
    CartHandler       *handlers.HttpCartHandler
    OrderHandler      *handlers.HttpOrderHandler
}

func NewContainer(db *gorm.DB) *Container {
    // Repositories
    userRepo := adapters.NewGormUserRepository(db)
    productRepo := adapters.NewGormProductRepository(db)
    categoriesRepo := adapters.NewGormCategoryRepository(db)
    cartRepo := adapters.NewGormCartRepository(db)
    orderRepo := adapters.NewGormOrderRepository(db)

    // Services
    passwordService := hash.NewPasswordService()
    userService := usecases.NewUserService(userRepo, passwordService)
    productService := usecases.NewProductService(productRepo)
    categoriesService := usecases.NewCategoryService(categoriesRepo)
    cartService := usecases.NewCartService(cartRepo,productRepo,orderRepo)
    orderService := usecases.NewOrderService(orderRepo)

    // Handlers
    return &Container{
        UserHandler:       handlers.NewHttpUserHandler(userService),
        ProductHandler:    handlers.NewHttpProductHandler(productService),
        CategoriesHandler: handlers.NewHttpCategoryHandler(categoriesService),
        CartHandler:       handlers.NewHttpCartHandler(cartService),
        OrderHandler:      handlers.NewHttpOrderHandler(orderService),
        // HealthHandler:     adapters.NewHealthHandler(db),
    }
}