package dto

type GetExpenseRes struct {
	Expenses []Expense `json:"expenses"`
	Total int64 `json:"total"`
}