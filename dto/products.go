package dto

type Products struct {
	Id    uint   `json:"id"`
	Price uint   `json:"price"`
	Item  string `json:"item"`
}