package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/scientist-v08/crackers/constants"
	"github.com/scientist-v08/crackers/controller"
	"github.com/scientist-v08/crackers/middleware"
)

func RegisterPriceList(app *fiber.App) {
	app.Get("/get/price-list", middleware.RequireAnyRole(constants.RoleUser, constants.RoleAdmin), controller.GetAllPrices)
}