package api

import (
	"github.com/Fal2o/E-Commerce_API/infastructure/config"
	"github.com/Fal2o/E-Commerce_API/infastructure/container"
	"github.com/Fal2o/E-Commerce_API/infastructure/routes"
	"github.com/Fal2o/E-Commerce_API/infastructure/server"
	database "github.com/Fal2o/E-Commerce_API/migrations"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load config
	cfg := config.LoadConfig()
	log.Printf("🚀 Starting E-Commerce API in %s mode", cfg.App.Environment)

	// Init database
	db := database.InitDB()
	database.AutoMigrate(db)

	// Init dependencies
	c := container.NewContainer(db)

	// Create server
	app := server.NewFiberApp(cfg)

	// Setup routes
	routes.Setup(app, c, cfg)

	// Start server
	go func() {
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			log.Fatalf("❌ Server failed: %v", err)
		}
	}()

	log.Printf("✅ Server running on port %s", cfg.Server.Port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down...")
	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Fatalf("❌ Forced shutdown: %v", err)
	}
	log.Println("✅ Shutdown complete")
}
