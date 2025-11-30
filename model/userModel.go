package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email string `gorm:"unique"`
	Password string
	Roles []string `gorm:"serializer:json"`
}