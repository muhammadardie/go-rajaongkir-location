package handlers

import (
	"fmt"
	"go-rajaongkir-location/config"
	"go-rajaongkir-location/dto"
	"go-rajaongkir-location/models"
	"go-rajaongkir-location/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllProvinces(c *gin.Context) {
	var provinces []models.Province
	var result []dto.ProvinceResponse
	id := c.Query("id")

	db := config.DB

	if id != "" {
		db = db.Where("province_id = ?", id)
	}

	if err := db.Find(&provinces).Error; err != nil {
		response.ErrorResponse(c, "Error retrieving provinces", http.StatusInternalServerError)
		return
	}

	for _, province := range provinces {
		result = append(result, dto.ProvinceResponse{
			ProvinceID:   fmt.Sprintf("%d", province.ProvinceID),
			ProvinceName: province.ProvinceName,
		})
	}

	if id != "" && len(result) == 1 {
		response.SuccessResponse(c, result[0])
		return
	}

	response.SuccessResponse(c, result)
}
