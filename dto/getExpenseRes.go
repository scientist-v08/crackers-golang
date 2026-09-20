package dto

import "github.com/scientist-v08/crackers/model"

type GetExpenseRes struct {
	Expenses []model.Expense
	Total int64
}