package model

import (
	"gorm.io/gorm"
)

type Billing struct {
	gorm.Model
	User         	string
	Mobile			string
	GrandTotal      int32
}