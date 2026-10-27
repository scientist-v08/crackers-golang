package dto

type Item struct {
	SlNo     int     `json:"slNo"`
	Item     string  `json:"item" binding:"required"`
	MRPOrNet float64 `json:"mrpOrNet" binding:"required,gt=0"`
	Quantity float64 `json:"quantity" binding:"required,gt=0"`
	Discount string  `json:"discount"`
	SubTotal float64 `json:"subTotal" binding:"required,gt=0"`
}