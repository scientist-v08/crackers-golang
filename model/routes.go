package model

type Routes struct {
	Id      uint   `gorm:"primaryKey;autoIncrement"`
	Route   string `gorm:"not null"`
	Heading string `gorm:"not null"`
	Role    string `gorm:"not null"`
}