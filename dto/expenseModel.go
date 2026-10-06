package dto

type Expense struct {
	ID               int64  `json:"id"`
	ReasonForExpense string `json:"reasonForExpense"`
	Amount           int32  `json:"amount"`
}