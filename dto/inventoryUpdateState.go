package dto

import "github.com/scientist-v08/crackers/constants"

type UpdateInventoryState struct {
	State 			constants.InventoryState
	SubTotal        int32
}