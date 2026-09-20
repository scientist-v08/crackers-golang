package model

type Billing struct {
	ID           uint64
	User         string
	Mobile       string
	GrandTotal   int32
	FinalizedAmt int32
	Purchases    []Purchases
}