package model

import "github.com/scientist-v08/crackers/constants"

type UpdateInventoryRequest struct {
	ID             uint                     `json:"ID" binding:"required"`
	BrandOrCompany string                   `json:"BrandOrCompany" binding:"required"`
	Item           string                   `json:"Item" binding:"required"`
	NumOfBoxes     int32                    `json:"NumOfBoxes"`
	NumOfCartons   int32                    `json:"NumOfCartons" binding:"required"`
	PricePerCarton int32                    `json:"PricePerCarton" binding:"required"`
	State          constants.InventoryState `json:"State" binding:"required"`
	// SubTotal is NOT accepted from frontend — we calculate it
}