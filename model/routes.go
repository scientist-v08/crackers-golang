package model

type Routes struct {
	Id      uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Route   string `gorm:"not null" json:"route"`
	Heading string `gorm:"not null" json:"heading"`
	Role    string `gorm:"not null" json:"role"`
}