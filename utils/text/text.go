package text

import "strings"

// parseCityName separates the 'Kota' or 'Kabupaten' prefix from the city_name
func ParseCityName(cityName string) (string, string) {
	var cityType, cityNameResult string

	if strings.HasPrefix(cityName, "Kota ") {
		cityType = "Kota"
		cityNameResult = strings.TrimPrefix(cityName, "Kota ")
	} else if strings.HasPrefix(cityName, "Kabupaten ") {
		cityType = "Kabupaten"
		cityNameResult = strings.TrimPrefix(cityName, "Kabupaten ")
	} else {
		cityType = ""
		cityNameResult = cityName
	}

	return cityType, cityNameResult
}
