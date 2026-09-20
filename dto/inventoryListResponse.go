package dto

import "github.com/scientist-v08/crackers/model"

type InventoryListResponse struct {
	InventoryItems []model.Inventory `json:"inventoryItems"`
	Total          int64             `json:"total"` // sum of sub_total
	TotalElements  int64             `json:"totalElements"`
}