package handlers

import (
	"fmt"
	"go-rajaongkir-location/config"
	"go-rajaongkir-location/dto"
	"go-rajaongkir-location/models"
	"go-rajaongkir-location/utils/response"
	"go-rajaongkir-location/utils/text"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllCity(c *gin.Context) {
	var cities []models.City
	var result []dto.CityResponse
	provinceID := c.Query("province")
	id := c.Query("id")

	db := config.DB.Preload("Province").Order("city_id ASC")

	if provinceID != "" {
		db = db.Where("province_id = ?", provinceID)
	}

	if id != "" {
		db = db.Where("city_id = ?", id)
	}

	if err := db.Find(&cities).Error; err != nil {
		response.ErrorResponse(c, "Failed to fetch cities", http.StatusInternalServerError)
		return
	}

	if len(cities) == 0 {
		response.SuccessResponse(c, []dto.CityResponse{})
		return
	}

	for _, city := range cities {
		// Use the helper function to separate the 'Kota' or 'Kabupaten' prefix
		cityType, cityName := text.ParseCityName(city.CityName)

		result = append(result, dto.CityResponse{
			CityID:     fmt.Sprintf("%d", city.CityID),
			ProvinceID: fmt.Sprintf("%d", city.ProvinceID),
			Province:   city.Province.ProvinceName,
			Type:       cityType,
			CityName:   cityName,
			PostalCode: city.PostalCode,
		})
	}

	response.SuccessResponse(c, result)
}
