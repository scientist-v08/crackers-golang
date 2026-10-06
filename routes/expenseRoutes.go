package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/scientist-v08/crackers/constants"
	"github.com/scientist-v08/crackers/controller"
	"github.com/scientist-v08/crackers/middleware"
)

func RegisterExpenseRoutes(r *fiber.App) {
	r.Post("/post/expenses", middleware.RequireAnyRole(constants.RoleAdmin), controller.CreateExpense)
	r.Get("/get/expenses", middleware.RequireAnyRole(constants.RoleAdmin), controller.GetAllExpenses)
}