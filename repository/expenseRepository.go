package repository

import (
	"github.com/scientist-v08/crackers/dto"
	"github.com/scientist-v08/crackers/initializers"
	"github.com/scientist-v08/crackers/model"
)

func CreateExpenses(exp *dto.ExpenseReqBody) error {
	expense := model.Expense{
		ReasonForExpense: exp.ReasonForExpense,
		Amount:           exp.Amount,
	}

	if err := initializers.DB.Create(&expense).Error; err != nil {
		return err
	}

	return nil
}

func GetAllExpenses() ([]model.Expense, error) {
	var resOfAllExpenses []model.Expense
	if err := initializers.DB.Find(&resOfAllExpenses).Error; err != nil {
		return nil, err
	}
	return resOfAllExpenses, nil
}

func GetAllExpensesTotalAmt() (int64, error) {
	var totalAmt int64
	if err := initializers.DB.Model(&model.Expense{}).Select("COALESCE(SUM(amount), 0)").Scan(&totalAmt).Error; err != nil {
		return 0, err
	}
	return totalAmt, nil
}