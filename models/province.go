package models

type Province struct {
	ProvinceID   int    `gorm:"column:province_id;primaryKey"`
	ProvinceName string `gorm:"column:province_name"`
}

func (Province) TableName() string {
	return "provinces"
}
