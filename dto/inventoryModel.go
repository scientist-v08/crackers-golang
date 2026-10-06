package dto

import (
	"github.com/scientist-v08/crackers/constants"
)

type Inventory struct {
	ID              uint64
	BrandOrCompany  string
	State           constants.InventoryState
	Item            string
	NumOfBoxes		int32
	NumOfCartons    int32
	PricePerCarton  int32
	SubTotal        int32
}