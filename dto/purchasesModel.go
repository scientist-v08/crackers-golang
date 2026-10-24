package dto

type Purchases struct {
	ID        uint64
	BillingID uint64
	Billing   Billing
	Mobile    string
	MrpOrNet  int32
	Item      string
	Quantity  int32
	Discount  string
	SubTotal  int32
}