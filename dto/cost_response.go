package dto

// CostRequest represents the incoming v1 API request
type CostRequest struct {
	Origin          string `form:"origin" binding:"required"`
	Destination     string `form:"destination" binding:"required"`
	Weight          string `form:"weight" binding:"required"`
	Courier         string `form:"courier" binding:"required"`
	OriginType      string `form:"originType"`
	DestinationType string `form:"destinationType"`
}

// V2Response represents the v2 API response format
type V2Response struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"meta"`
	Data []struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Service     string `json:"service"`
		Description string `json:"description"`
		Cost        int    `json:"cost"`
		ETD         string `json:"etd"`
	} `json:"data"`
}

type V2CostData struct {
	Code        string `json:"code"`
	Cost        int    `json:"cost"`
	Description string `json:"description"`
	ETD         string `json:"etd"`
	Name        string `json:"name"`
	Service     string `json:"service"`
}

type V1Response struct {
	Code  string   `json:"code"`
	Name  string   `json:"name"`
	Costs []V1Cost `json:"costs"`
}

type V1Cost struct {
	Service     string         `json:"service"`
	Description string         `json:"description"`
	Cost        []V1CostDetail `json:"cost"`
}

type V1CostDetail struct {
	Value int    `json:"value"`
	ETD   string `json:"etd"`
	Note  string `json:"note"`
}
