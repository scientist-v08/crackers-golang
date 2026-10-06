package service

import (
	"fmt"
	"strings"
	"sync"

	"github.com/scientist-v08/crackers/db"
	"github.com/scientist-v08/crackers/dto"
	"github.com/scientist-v08/crackers/repository"
)

var (
	ErrReasonTooShort = "reasonForExpense must be at least 5 characters long after trimming spaces"
	ErrAmountTooLow   = "amount must be at least 10"
)

func validateExpenseRequest(req *dto.ExpenseReqBody) error {

	// Validate ReasonForExpense: trim and check minimum length 5
	trimmedReason := strings.TrimSpace(req.ReasonForExpense)
	if len(trimmedReason) < 5 {
		return fmt.Errorf("%s", ErrReasonTooShort)
	}

	// Validate Amount: must be at least 10
	if req.Amount < 10 {
		return fmt.Errorf("%s", ErrAmountTooLow)
	}

	return nil
}

func CreateExpenseService(req dto.ExpenseReqBody) error {
	if err := validateExpenseRequest(&req); err != nil {
		return err
	}

	if err := repository.CreateExpenses(&req); err != nil {
		return err
	}

	return nil
}

func GetAllExpenseService() (dto.GetExpenseRes, error) {
	// 1. Create variables for obtaining errors and for go routines
	var allExpensesErr, totalErr error
	var wg sync.WaitGroup
	var resOfAllExpenses []db.Expense
	var totalAmount int64

	wg.Add(2)

	// 2. Obtain all the results
	go func() {
		defer wg.Done()
		resOfAllExpenses, allExpensesErr = repository.GetAllExpenses()
	}()
	
	// 3. Finding the total budget
	go func() {
		defer wg.Done()
		totalAmount, totalErr = repository.GetAllExpensesTotalAmt()
	}()

	// 4. Wait until both the go routines finish.
	wg.Wait()
	
	// 5. Return any errors
	if allExpensesErr != nil {
		return dto.GetExpenseRes{}, allExpensesErr
	}
	if totalErr != nil {
		return dto.GetExpenseRes{}, totalErr
	}

	// 6. Map db.Expense → dto.Expense
	expenses := make([]dto.Expense, 0, len(resOfAllExpenses))
	for _, e := range resOfAllExpenses {
		expenses = append(expenses, dto.Expense{
			ID:               e.ID,                 // adjust field names if different
			ReasonForExpense: e.ReasonForExpense.String,   // adjust if the db field is named differently
			Amount:           e.Amount.Int32,             // cast if needed (int32 vs int64 etc.)
		})
	}

	// 7. Return results
	return dto.GetExpenseRes{Expenses: expenses, Total: totalAmount}, nil
}