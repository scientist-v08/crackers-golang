package dto

import "github.com/scientist-v08/crackers/constants"

type UpdateInventoryRequest struct {
	ID             uint                     `json:"id" binding:"required"`
	BrandOrCompany string                   `json:"brandOrCompany" binding:"required"`
	Item           string                   `json:"item" binding:"required"`
	NumOfBoxes     int32                    `json:"numOfBoxes"`
	NumOfCartons   int32                    `json:"numOfCartons" binding:"required"`
	PricePerCarton int32                    `json:"pricePerCarton" binding:"required"`
	State          constants.InventoryState `json:"state" binding:"required"`
	// SubTotal is NOT accepted from frontend — we calculate it
}