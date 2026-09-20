package dto

type ExpenseReqBody struct {
	ReasonForExpense string `json:"reasonForExpense"`
	Amount           int32  `json:"amount"`
}