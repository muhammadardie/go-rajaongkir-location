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

	response.SuccessResponse(c, result)
}

func GetProvinceByID(c *gin.Context) {
	id := c.Param("id")

	var province models.Province

	if err := config.DB.Where("province_id = ?", id).First(&province).Error; err != nil {
		response.ErrorResponse(c, "Province not found", http.StatusNotFound)
		return
	}

	result := dto.ProvinceResponse{
		ProvinceID:   fmt.Sprintf("%d", province.ProvinceID),
		ProvinceName: province.ProvinceName,
	}

	response.SuccessResponse(c, result)
}
