package model

type Purchases struct {
	ID        uint64
	BillingID uint64  `gorm:"not null;index"`
	Billing   Billing `gorm:"foreignKey:BillingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Mobile    string  `gorm:"index"`
	MrpOrNet  int32
	Item      string
	Quantity  int32
	Discount  string
	SubTotal  int32
}