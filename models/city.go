package models

type City struct {
	CityID     int      `gorm:"column:city_id;primaryKey" json:"city_id"`
	ProvinceID int      `gorm:"column:province_id" json:"province_id"`
	CityName   string   `gorm:"column:city_name" json:"city_name"`
	PostalCode string   `gorm:"column:postal_code" json:"postal_code"`
	Province   Province `gorm:"foreignKey:ProvinceID;references:ProvinceID" json:"-"`
}

func (City) TableName() string {
	return "cities"
}
