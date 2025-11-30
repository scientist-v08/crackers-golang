package model

import "gorm.io/gorm"

type Expense struct {
	gorm.Model
	ReasonForExpense string
	Amount			 int32
}