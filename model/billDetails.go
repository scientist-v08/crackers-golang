package model

type BillDetails struct {
	User         string  `json:"user" binding:"required"`
	Mobile       string  `json:"mobile" binding:"required"`
	GrandTotal   float64 `json:"grandTotal" binding:"required,gt=0"`
	FinalizedAmt int32   `json:"finalizedAmt"`
	Discount     string  `json:"discount" binding:"required"`
	BillItems    []Item  `json:"billItems" binding:"required,min=1"`
}