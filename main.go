package main

import (
	"go-rajaongkir-location/config"
	"go-rajaongkir-location/middleware"
	"go-rajaongkir-location/routes"
	"log"

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
	router.SetTrustedProxies(nil)
	routes.SetupRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
