package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/scientist-v08/crackers/db"
	"github.com/scientist-v08/crackers/dto"
	"github.com/scientist-v08/crackers/initializers"
)

func CreateExpenses(exp *dto.ExpenseReqBody) error {
	_, err := initializers.Queries.CreateExpense(context.Background(), db.CreateExpenseParams{
		ReasonForExpense: pgtype.Text{
			String: exp.ReasonForExpense,
			Valid:  true,
		},
		Amount: pgtype.Int4{
			Int32: exp.Amount,
			Valid:  true,
		},
	})

	return err
}

func GetAllExpenses() ([]db.Expense, error) {
	expenses, err := initializers.Queries.GetAllExpenses(context.Background())
	if err != nil {
		return nil, err
	}

	return expenses, nil
}

func GetAllExpensesTotalAmt() (int64, error) {
	totalAmt, err := initializers.Queries.GetAllExpensesTotalAmt(context.Background())
	if err != nil {
		return 0, err
	}

	return totalAmt, nil
}