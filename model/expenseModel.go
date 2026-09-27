package model

type Expense struct {
	ID               uint64 `json:"id"`
	ReasonForExpense string `json:"reasonForExpense"`
	Amount           int32  `json:"amount"`
}