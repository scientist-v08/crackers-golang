package dto

import "github.com/scientist-v08/crackers/constants"

type InventoryReqBody struct {
	BrandOrCompany string                   `json:"brandOrCompany"`
	State          constants.InventoryState `json:"state"`
	Item           string                   `json:"item"`
	NumOfBoxes     int32                    `json:"numOfBoxes" binding:"required,gt=0"`
	NumOfCartons   int32                    `json:"numOfCartons" binding:"required,gt=0"`
	PricePerCarton int32                    `json:"pricePerCarton" binding:"required,gt=0"`
	SubTotal       int32                    `json:"subtotal" binding:"required,gt=0"`
}