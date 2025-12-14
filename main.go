package main

import (
	"workmgmt-api/internal/config"
	"workmgmt-api/internal/handlers"
	"workmgmt-api/internal/router"
	"log"

	_ "workmgmt-api/docs" // swagger docs
)

// @title Workmgmt Api
// @version 1.0.0
// @description Work management system integrations
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8089
// @BasePath /api/v1

// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the access token.

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	apiHandler := handlers.NewAPIHandler(cfg)

	// Setup routes
	appRouter := router.NewRouter()
	router.RegisterHealthRoutes(appRouter, healthHandler)
	router.RegisterAPIRoutes(appRouter, apiHandler)
	router.RegisterSwaggerRoutes(appRouter)

	// Start server
	server := router.NewServer(cfg, appRouter)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
