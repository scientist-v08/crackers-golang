package dto

type TempBillItem struct {
	MrpOrNet int32   `json:"mrpOrNet"`
	Item     string  `json:"item"`
	Quantity int32   `json:"quantity"`
	Discount string  `json:"discount"`
	SubTotal float64 `json:"subTotal"`
}

type TempBillDetails struct {
	User          string         `json:"user" binding:"required"`
	Mobile        string         `json:"mobile" binding:"required"`
	GrandTotal    float64        `json:"grandTotal" binding:"required,gt=0"`
	FinalizedAmt  float64        `json:"finalizedAmt"`
	Comments      string         `json:"comments" binding:"min=5"`
	Discount      string         `json:"discount"` // top-level, currently unused in DB
	TempBillItems []TempBillItem `json:"tempBillItems" binding:"required,min=1"`
}

type TempBillIDRequest struct {
	ID int16 `json:"id" binding:"required"`
}

// Response for /all/get-saved
type TempBillResponse struct {
	ID           int16          `json:"id"`
	User         string         `json:"user"`
	Mobile       string         `json:"mobile"`
	Comments     string         `json:"comments"`
	GrandTotal   float64        `json:"grandTotal"`
	FinalizedAmt float64        `json:"finalizedAmt"`
	Items        []TempBillItem `json:"items"`
}

type CommentsAndFinalizedAmt struct {
	ID           int16   `json:"id"`
	Comments     string  `json:"comments"`
	FinalizedAmt float64 `json:"finalizedAmt"`
}