package middleware

import (
	"go-rajaongkir-location/utils/response"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	requiredAPIKey := os.Getenv("API_KEY")

	return func(c *gin.Context) {
		// Check Authorization header for new API flow
		authHeader := c.GetHeader("Authorization")
		rajaongkirKey := c.GetHeader("rajaongkir-key")

		// Legacy support: if only 'key' is present, allow
		if authHeader == "" && c.GetHeader("key") != "" {
			// Treat 'key' as RajaOngkir key
			c.Set("rajaongkir_key", c.GetHeader("key"))
			c.Next()
			return
		}

		// New auth required
		if requiredAPIKey == "" {
			// If no API_KEY configured, skip auth
			c.Next()
			return
		}

		// Authorization must be "Bearer <token>"
		const bearerPrefix = "Bearer "
		if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
			response.ErrorResponse(c, "Missing or invalid Authorization header", http.StatusUnauthorized)
			c.Abort()
			return
		}

		token := authHeader[len(bearerPrefix):]

		// Validate token (replace this with your actual logic)
		if token != requiredAPIKey {
			response.ErrorResponse(c, "Unauthorized: Invalid token", http.StatusUnauthorized)
			c.Abort()
			return
		}

		// Must also have RajaOngkir key in new flow
		if rajaongkirKey == "" {
			response.ErrorResponse(c, "Missing rajaongkir-key header", http.StatusBadRequest)
			c.Abort()
			return
		}

		c.Set("rajaongkir_key", rajaongkirKey)
		c.Next()
	}
}
