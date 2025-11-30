package model

import (
	"gorm.io/gorm"
)

type Purchases struct {
	gorm.Model
	BillingID 		uint   `gorm:"not null;index"`
	Mobile			string `gorm:"index"`
	MrpOrNet        int32
	Item			string
	Quantity		int32
	Discount		string
	SubTotal		int32
}