package dto

type SubdistrictResponse struct {
	SubdistrictID   string `json:"subdistrict_id"`
	ProvinceID      string `json:"province_id"`
	Province        string `json:"province"`
	CityID          string `json:"city_id"`
	City            string `json:"city"`
	Type            string `json:"type"`
	SubdistrictName string `json:"subdistrict_name"`
	PostalCode      string `json:"postal_code"`
}
