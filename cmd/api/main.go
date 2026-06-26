package main

import (
	"bast-request/internal/config"
	"bast-request/internal/routes"

	"github.com/gin-gonic/gin"
)

// @title BAST Request API
// @version 1.0
// @description This is the API server for BAST Request System.
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Initialize Database Connection
	config.ConnectDB()
	
	// Run Auto Migration
	config.AutoMigrate()

	// Seed dummy data
	config.SeedDB(config.DB)

	// Initialize Gin router
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Setup API Routes
	routes.SetupRoutes(r, config.DB)

	// Run the server
	r.Run(":8080")
}
