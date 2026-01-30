package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Fal2o/E-Commerce_API/docs"
	"github.com/Fal2o/E-Commerce_API/infrastructure/config"
	"github.com/Fal2o/E-Commerce_API/infrastructure/container"
	"github.com/Fal2o/E-Commerce_API/infrastructure/routes"
	"github.com/Fal2o/E-Commerce_API/infrastructure/server"
	database "github.com/Fal2o/E-Commerce_API/migrations" 
	"github.com/gofiber/swagger"
)

// @title E-Commerce API
// @version 1.0.0
// @description Production-ready RESTful API for e-commerce platform built with Go and Clean Architecture.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@ecommerce-api.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8000
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load config
	cfg := config.LoadConfig()
	log.Printf("Starting E-Commerce API in %s mode", cfg.App.Environment)

	// Define flags
	seedFlag := flag.Bool("seed", false, "Run database seeding")
	flag.Parse()

	// Init database
	db := database.InitDB()
	database.AutoMigrate(db)

	// Run seed if flag is provided
	if *seedFlag {
		log.Println("Running database seed...")
		database.SeedData(db)
		log.Println("Seed completed!")
		return
	}
	// Init dependencies
	c := container.NewContainer(db)

	// Create server
	app := server.NewFiberApp(cfg)
	app.Get("/swagger/*", swagger.HandlerDefault)
	// Setup routes
	routes.Setup(app, c, cfg)

	// Start server
	go func() {
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()
	log.Printf("Server running on port %s", cfg.Server.Port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Fatalf("Forced shutdown: %v", err)
	}
	log.Println("Shutdown complete")
}
