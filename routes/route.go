package routes

import (
	"go-rajaongkir-location/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/province/:id", handlers.GetProvinceByID)
	router.GET("/province", handlers.GetAllProvinces)
	router.GET("/city", handlers.GetAllCity)
	router.GET("/subdistrict", handlers.GetAllSubdistrict)
}
