package dto

type InventoryListResponse struct {
	InventoryItems []Inventory `json:"inventoryItems"`
	Total          int64       `json:"total"` // sum of sub_total
	TotalElements  int64       `json:"totalElements"`
}