package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/scientist-v08/crackers/constants"
	"github.com/scientist-v08/crackers/controller"
	"github.com/scientist-v08/crackers/middleware"
)

func RegisterExpenseRoutes(r *gin.Engine) {
	r.POST("/post/expenses", middleware.RequireAnyRole(constants.RoleAdmin), controller.CreateExpense)
	r.GET("/get/expenses", middleware.RequireAnyRole(constants.RoleAdmin), controller.GetAllExpenses)
}