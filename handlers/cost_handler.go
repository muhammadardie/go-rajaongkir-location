package handlers

import (
	"encoding/json"
	"fmt"
	"go-rajaongkir-location/config"
	"go-rajaongkir-location/dto"
	"go-rajaongkir-location/models"
	"go-rajaongkir-location/utils/response"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CourierMapping maps v1 courier names to v2 courier names
var CourierMapping = map[string]string{
	"jne":      "jne",
	"sicepat":  "sicepat",
	"ide":      "ide",
	"sap":      "sap",
	"jnt":      "jnt",
	"ninja":    "ninja",
	"tiki":     "tiki",
	"lion":     "lion",
	"anteraja": "anteraja",
	"pos":      "pos",
	"ncs":      "ncs",
	"rex":      "rex",
	"rpx":      "rpx",
	"sentral":  "sentral",
	"star":     "star",
	"wahana":   "wahana",
	"dse":      "dse",
}

// AllowedCouriers list of allowed couriers
var AllowedCouriers = []string{"jne", "sicepat", "ide", "sap", "jnt", "ninja", "tiki", "lion", "anteraja", "pos", "ncs", "rex", "rpx", "sentral", "star", "wahana", "dse"}

// GetCost handles the cost calculation endpoint
func GetCost(c *gin.Context) {
	var req dto.CostRequest

	rawKey, exists := c.Get("rajaongkir_key")
	if !exists {
		response.ErrorResponse(c, "Missing RajaOngkir API key", http.StatusBadRequest)
		return
	}

	apiKey, ok := rawKey.(string)
	if !ok || apiKey == "" {
		response.ErrorResponse(c, "Invalid RajaOngkir API key", http.StatusBadRequest)
		return
	}

	// Bind form data
	if err := c.ShouldBind(&req); err != nil {
		response.ErrorResponse(c, "Invalid request parameters", http.StatusBadRequest)
		return
	}

	// Validate courier(s) - can be single or multiple separated by ':'
	if !isValidCouriers(req.Courier) {
		response.ErrorResponse(c, "Invalid courier(s). Allowed couriers: "+strings.Join(AllowedCouriers, ", "), http.StatusBadRequest)
		return
	}

	// Validate weight
	if req.Weight <= 0 {
		response.ErrorResponse(c, "Invalid weight. Must be a positive integer", http.StatusBadRequest)
		return
	}

	// Get postal codes from subdistrict IDs
	originPostalCode, err := getPostalCodeBySubdistrictID(req.Origin)
	if err != nil {
		response.ErrorResponse(c, "Invalid origin subdistrict ID", http.StatusBadRequest)
		return
	}

	destinationPostalCode, err := getPostalCodeBySubdistrictID(req.Destination)
	if err != nil {
		response.ErrorResponse(c, "Invalid destination subdistrict ID", http.StatusBadRequest)
		return
	}

	// Map courier name for v2 API
	v2Courier := mapCouriersForV2(req.Courier)
	weight := strconv.Itoa(req.Weight)

	// Call RajaOngkir API v2
	v2Response, err := callRajaOngkirV2(apiKey, originPostalCode, destinationPostalCode, weight, v2Courier)
	if err != nil {
		log.Printf(`{"level":"error","msg":"callRajaOngkirV2 failed","error":"%v"}`, err)
		response.ErrorResponse(c, "Failed to get cost data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Transform v2 response to v1 format
	v1Response := transformV2ToV1(v2Response)

	// Return the response in v1 format
	response.SuccessResponse(c, v1Response)
}

// isValidCouriers checks if all couriers in the string are valid (supports single or multiple separated by ':')
func isValidCouriers(couriers string) bool {
	courierList := strings.Split(couriers, ":")

	for _, courier := range courierList {
		courier = strings.TrimSpace(courier)
		if courier == "" {
			continue
		}

		found := false
		for _, allowedCourier := range AllowedCouriers {
			if courier == allowedCourier {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

// mapCouriersForV2 maps courier names for v2 API (keeps the same format as v1)
func mapCouriersForV2(couriers string) string {
	courierList := strings.Split(couriers, ":")
	var mappedCouriers []string

	for _, courier := range courierList {
		courier = strings.TrimSpace(courier)
		if courier == "" {
			continue
		}

		if mappedCourier, exists := CourierMapping[courier]; exists {
			mappedCouriers = append(mappedCouriers, mappedCourier)
		}
	}

	return strings.Join(mappedCouriers, ":")
}

// getPostalCodeBySubdistrictID retrieves postal code from database using subdistrict ID
func getPostalCodeBySubdistrictID(subdistrictID string) (string, error) {
	var subdistrict models.Subdistrict

	// Convert string ID to int
	id, err := strconv.Atoi(subdistrictID)
	if err != nil {
		return "", fmt.Errorf("invalid subdistrict ID format")
	}

	// Query database with preload to get city data
	if err := config.DB.Where("subdistrict_id = ?", id).First(&subdistrict).Error; err != nil {
		return "", fmt.Errorf("subdistrict not found")
	}

	// Return postal code from V2CostData
	return subdistrict.PostalCode, nil
}

// callRajaOngkirV2 makes the actual API call to RajaOngkir v2
func callRajaOngkirV2(apiKey, origin, destination, weight, courier string) (*dto.V2Response, error) {
	RajaOngkirV2CostURL := os.Getenv("RAJAONGKIR_V2_COST_URL")

	// Prepare form data
	data := url.Values{}
	data.Set("origin", origin)
	data.Set("destination", destination)
	data.Set("weight", weight)
	data.Set("courier", courier)

	// Create HTTP request
	req, err := http.NewRequest("POST", RajaOngkirV2CostURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("key", apiKey)

	// Make HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Check if request was successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	var v2Response dto.V2Response
	if err := json.Unmarshal(body, &v2Response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	// Check if API response indicates success
	if v2Response.Meta.Code != 200 {
		return nil, fmt.Errorf("API error: %s (code %d)", v2Response.Meta.Message, v2Response.Meta.Code)
	}

	return &v2Response, nil
}

// transformV2ToV1 converts v2 API response format to v1 format
func transformV2ToV1(v2Response *dto.V2Response) []dto.V1Response {
	mapResponse := make(map[string]*dto.V1Response)

	// Group costs by courier code
	for _, costData := range v2Response.Data {
		if courierResult, exists := mapResponse[costData.Code]; exists {
			// Add cost to existing courier
			courierResult.Costs = append(courierResult.Costs, dto.V1Cost{
				Service:     costData.Service,
				Description: costData.Description,
				Cost: []dto.V1CostDetail{
					{
						Value: costData.Cost,
						ETD:   costData.ETD,
						Note:  "",
					},
				},
			})
		} else {
			// Create new courier entry
			mapResponse[costData.Code] = &dto.V1Response{
				Code: costData.Code,
				Name: costData.Name,
				Costs: []dto.V1Cost{
					{
						Service:     costData.Service,
						Description: costData.Description,
						Cost: []dto.V1CostDetail{
							{
								Value: costData.Cost,
								ETD:   costData.ETD,
								Note:  "",
							},
						},
					},
				},
			}
		}
	}

	// Convert map to slice
	var results []dto.V1Response
	for _, courierResult := range mapResponse {
		results = append(results, *courierResult)
	}

	return results
}
