package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Status struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type Response struct {
	Status  Status      `json:"status"`
	Results interface{} `json:"results"`
}

type RajaOngkirResponse struct {
	RajaOngkir Response `json:"rajaongkir"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
	var results interface{}
	switch v := data.(type) {
	case []interface{}:
		results = v // Multiple results
	default:
		results = data // Single result
	}

	response := RajaOngkirResponse{
		RajaOngkir: Response{
			Status: Status{
				Code:        http.StatusOK,
				Description: "OK",
			},
			Results: results,
		},
	}

	c.JSON(http.StatusOK, response)
}

func ErrorResponse(c *gin.Context, message string, statusCode int) {
	response := RajaOngkirResponse{
		RajaOngkir: Response{
			Status: Status{
				Code:        statusCode,
				Description: message,
			},
			Results: nil,
		},
	}

	c.JSON(statusCode, response)
}
