package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scientist-v08/crackers/initializers"
	"github.com/scientist-v08/crackers/model"
)

// ---------- Frontend interfaces mirrored in Go --------------
type ExpenseReqBody struct {
	ReasonForExpense string `json:"reasonForExpense"`
	Amount           int32  `json:"amount"`
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

	// 2. Map to actual Expense model
	expense := model.Expense{
		ReasonForExpense: req.ReasonForExpense,
		Amount:           req.Amount,
	}

	// 3. Save the data in the DB
	if err := initializers.DB.Create(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the data in the DB"})
		return
	}

	// 4. Send success message
	c.JSON(http.StatusCreated, gin.H{"Success": "Item saved to inventory"})
}

func GetAllExpenses(c *gin.Context) {

	// 1. Obtain all the results
	var resOfAllExpenses []model.Expense
	if err := initializers.DB.Find(&resOfAllExpenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed To fetch expenses",
		})
		return
	}

	// 2. Finding the total budget
	var totalAmount int64
	if err := initializers.DB.Model(&model.Expense{}).Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total"})
        return
	}

	// 3. Return all the results
	c.JSON(http.StatusOK, gin.H{
		"expenses": resOfAllExpenses,
		"total": totalAmount,
	})
}
