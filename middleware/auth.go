package middleware

import (
	"go-rajaongkir-location/utils/response"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	requiredAPIKey := os.Getenv("API_KEY")

	// skip auth when env API_KEY not supplied
	if requiredAPIKey == "" {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")
		if apiKey == "" || apiKey != requiredAPIKey {
			response.ErrorResponse(c, "Invalid or missing API key", http.StatusUnauthorized)
			c.Abort()
			return
		}
		c.Next()
	}
}
