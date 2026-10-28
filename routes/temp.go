package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/scientist-v08/crackers/constants"
	"github.com/scientist-v08/crackers/controller"
	"github.com/scientist-v08/crackers/middleware"
)

func RegisterTempRoutes(app *fiber.App) {
	// or use a group if you prefer: api := app.Group("/all")
	app.Post("/all/save-temp", middleware.RequireAnyRole(constants.RoleUser, constants.RoleAdmin), controller.SaveTempBill)
	app.Post("/all/get-saved", middleware.RequireAnyRole(constants.RoleUser, constants.RoleAdmin), controller.GetSavedTempBill)
	app.Post("/all/delete-saved", middleware.RequireAnyRole(constants.RoleUser, constants.RoleAdmin), controller.DeleteSavedTempBill)
	app.Get("/all/comments-and-amounts", middleware.RequireAnyRole(constants.RoleUser, constants.RoleAdmin), controller.GetTempBillingCommentsAndFinalized)
}