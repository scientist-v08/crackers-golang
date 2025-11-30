package model

import (
	"github.com/scientist-v08/crackers/constants"
	"gorm.io/gorm"
)

type Inventory struct {
	gorm.Model
	BrandOrCompany  string
	State           constants.InventoryState `gorm:"type:varchar(20);index"`
	Item            string `gorm:"index"`
	NumOfBoxes		int32
	NumOfCartons    int32
	PricePerCarton  int32
	SubTotal        int32
}