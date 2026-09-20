package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scientist-v08/crackers/dto"
	"github.com/scientist-v08/crackers/service"
)

// ---------- Gin Handler -------------
func CreateExpense(c *gin.Context) {

	// 1. Obtain the request body
	var req dto.ExpenseReqBody
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect request body",
		})
		return
	}

	// 2. Save the data in the DB
	if err := service.CreateExpenseService(req); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	// 3. Send success message
	c.JSON(http.StatusCreated, gin.H{"Success": "Item saved to inventory"})
}

func GetAllExpenses(c *gin.Context) {
	// 1. Create variables
	var res dto.GetExpenseRes
	var err error

	// 2. Obtain results
	res, err = service.GetAllExpenseService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// 3. Return results
	c.JSON(http.StatusOK, res)
}
