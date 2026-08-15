package model

type PriceList struct {
	Id    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Price uint   `json:"price"`
	Item  string `json:"item"`
	Brand string `json:"brand"`
}