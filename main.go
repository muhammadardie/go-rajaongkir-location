package main

import (
	"go-rajaongkir-location/config"
	"go-rajaongkir-location/middleware"
	"go-rajaongkir-location/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}
	config.InitDB()

	rateLimiter := middleware.CreateRateLimiter()
	router := gin.Default()
	router.Use(rateLimiter.RateLimit())
	router.Use(middleware.UmamiAnalyticsMiddleware())
	router.SetTrustedProxies(nil)
	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000" // default port if not set
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
