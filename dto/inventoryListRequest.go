package dto

type InventoryListRequest struct {
	PageNumber int
	PageSize   int
	SortBy     string
	Order      string
	Title      string
	State      string
}