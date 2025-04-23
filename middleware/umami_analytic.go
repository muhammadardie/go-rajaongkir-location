package middleware

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func UmamiAnalyticsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("ENABLE_UMAMI") != "true" {
			c.Next()
			return
		}

		go func() {
			host := os.Getenv("UMAMI_HOST")
			if host == "" {
				return
			}

			payload := map[string]interface{}{
				"payload": map[string]interface{}{
					"website":  os.Getenv("UMAMI_WEBSITE_ID"),
					"hostname": os.Getenv("UMAMI_HOSTNAME"),
					"url":      c.Request.URL.Path,
					"referrer": c.Request.Referer(),
					"language": c.GetHeader("Accept-Language"),
					"screen":   "1920x1080",
					"title":    c.Request.URL.Path,
					"name":     "pageview",
				},
				"type": "event",
			}

			jsonBody, _ := json.Marshal(payload)

			req, err := http.NewRequest("POST", host+"/api/send", bytes.NewBuffer(jsonBody))
			if err != nil {
				log.Println("Failed to create request for Umami:", err)
				return
			}

			req.Header.Set("User-Agent", os.Getenv("UMAMI_USER_AGENT"))
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{Timeout: 2 * time.Second}
			client.Do(req)
		}()

		c.Next()
	}
}
