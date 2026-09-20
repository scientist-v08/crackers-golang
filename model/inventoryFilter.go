package model

import "github.com/scientist-v08/crackers/constants"

type InventoryFilter struct {
	Title string
	State constants.InventoryState
}