package main

import (
	"go-rajaongkir-location/config"
	"go-rajaongkir-location/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()

	router := gin.Default()
	router.SetTrustedProxies(nil)
	routes.SetupRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
