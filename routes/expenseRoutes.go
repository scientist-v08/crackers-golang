package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/scientist-v08/crackers/constants"
	"github.com/scientist-v08/crackers/controller"
	"github.com/scientist-v08/crackers/middleware"
)

func RegisterExpenseRoutes(r *fiber.App) {
	r.Post("/admin/expenses", middleware.RequireAnyRole(constants.RoleAdmin), controller.CreateExpense)
	r.Get("/admin/expenses", middleware.RequireAnyRole(constants.RoleAdmin), controller.GetAllExpenses)
}