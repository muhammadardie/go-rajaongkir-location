package models

type Subdistrict struct {
	SubdistrictID   int     `gorm:"column:subdistrict_id;primaryKey"`
	CityID          int     `gorm:"column:city_id"`
	SubdistrictName *string `gorm:"column:subdistrict_name"`
	City            City    `gorm:"foreignKey:CityID;references:CityID"`
	PostalCode      string  `gorm:"column:postal_code" json:"postal_code"`
}

func (Subdistrict) TableName() string {
	return "subdistricts"
}
