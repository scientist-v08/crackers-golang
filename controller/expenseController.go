package controller

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/scientist-v08/crackers/initializers"
	"github.com/scientist-v08/crackers/model"
)

// ---------- Frontend interfaces mirrored in Go --------------
type ExpenseReqBody struct {
	ReasonForExpense string `json:"reasonForExpense"`
	Amount           int32  `json:"amount"`
}

var (
	ErrReasonTooShort = "reasonForExpense must be at least 5 characters long after trimming spaces"
	ErrAmountTooLow   = "amount must be at least 10"
)

func validateExpenseRequest(req *ExpenseReqBody) error {

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

// ---------- Gin Handler -------------
func CreateExpense(c *gin.Context) {

	// 1. Obtain the request body
	var req ExpenseReqBody
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect request body",
		})
		return
	}

	// 2. Validate the data inputs.
	if err := validateExpenseRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 3. Map to actual Expense model
	expense := model.Expense{
		ReasonForExpense: req.ReasonForExpense,
		Amount:           req.Amount,
	}

	// 4. Save the data in the DB
	if err := initializers.DB.Create(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the data in the DB"})
		return
	}

	// 5. Send success message
	c.JSON(http.StatusCreated, gin.H{"Success": "Item saved to inventory"})
}

func GetAllExpenses(c *gin.Context) {

	// 1. Create variables for obtaining errors and for go routines
	var allExpensesErr, totalErr error
	var wg sync.WaitGroup
	var resOfAllExpenses []model.Expense
	var totalAmount int64

	wg.Add(2)

	// 2. Obtain all the results
	go func() {
		defer wg.Done()
		allExpensesErr = initializers.DB.Find(&resOfAllExpenses).Error
	}()
	
	// 3. Finding the total budget
	go func() {
		defer wg.Done()
		totalErr = initializers.DB.Model(&model.Expense{}).Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount).Error
	}()

	// 4. Wait until both the go routines finish.
	wg.Wait()
	
	// 5. Return all the results
	if allExpensesErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed To fetch expenses",
		})
		return
	}
	if totalErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"expenses": resOfAllExpenses,
		"total": totalAmount,
	})
}
