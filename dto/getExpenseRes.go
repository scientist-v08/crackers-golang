package dto

import "github.com/scientist-v08/crackers/model"

type GetExpenseRes struct {
	Expenses []model.Expense `json:"expenses"`
	Total int64 `json:"total"`
}