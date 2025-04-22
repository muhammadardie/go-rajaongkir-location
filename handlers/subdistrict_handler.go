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

func GetAllSubdistrict(c *gin.Context) {
	var subdistricts []models.Subdistrict
	var result []dto.SubdistrictResponse
	provinceID := c.Query("province")
	cityID := c.Query("city")
	id := c.Query("id")

	db := config.DB.Preload("City").Preload("City.Province").Order("subdistrict_id ASC")

	if provinceID != "" {
		db = db.Joins("JOIN cities ON subdistricts.city_id = cities.city_id").
			Where("cities.province_id = ?", provinceID)
	}

	if cityID != "" {
		db = db.Where("subdistricts.city_id = ?", cityID)
	}

	if id != "" {
		db = db.Where("subdistrict_id = ?", id)
	}

	if err := db.Find(&subdistricts).Error; err != nil {
		response.ErrorResponse(c, "Failed to fetch subdistricts", http.StatusInternalServerError)
		return
	}

	if len(subdistricts) == 0 {
		response.SuccessResponse(c, []dto.SubdistrictResponse{})
		return
	}

	for _, subdistrict := range subdistricts {
		provinceName := subdistrict.City.Province.ProvinceName
		cityType, cityName := text.ParseCityName(subdistrict.City.CityName)

		result = append(result, dto.SubdistrictResponse{
			SubdistrictID:   fmt.Sprintf("%d", subdistrict.SubdistrictID),
			ProvinceID:      fmt.Sprintf("%d", subdistrict.City.ProvinceID),
			Province:        provinceName,
			CityID:          fmt.Sprintf("%d", subdistrict.CityID),
			City:            cityName,
			Type:            cityType,
			SubdistrictName: *subdistrict.SubdistrictName,
			PostalCode:      subdistrict.City.PostalCode,
		})
	}

	if id != "" && len(result) == 1 {
		response.SuccessResponse(c, result[0])
		return
	}

	response.SuccessResponse(c, result)
}
