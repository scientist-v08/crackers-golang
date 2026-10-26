package dto

import (
	"github.com/scientist-v08/crackers/constants"
)

type Inventory struct {
	ID             uint64                   `json:"id"`
	BrandOrCompany string                   `json:"brandOrCompany"`
	State          constants.InventoryState `json:"state"`
	Item           string                   `json:"item"`
	NumOfBoxes     int32                    `json:"numOfBoxes"`
	NumOfCartons   int32                    `json:"numOfCartons"`
	PricePerCarton int32                    `json:"pricePerCarton"`
	SubTotal       int32                    `json:"subTotal"`
}