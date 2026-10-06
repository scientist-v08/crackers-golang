package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/scientist-v08/crackers/constants"
	"github.com/scientist-v08/crackers/controller"
	"github.com/scientist-v08/crackers/middleware"
)

func RegisterBillingRoutes(r *fiber.App) {
	r.Post("/billing", middleware.RequireAnyRole(constants.RoleUser, constants.RoleAdmin), controller.CreateBillHandler)
	r.Post("/billing/preview", middleware.RequireAnyRole(constants.RoleUser, constants.RoleAdmin), controller.CreatePreviewBillHandler)
}
