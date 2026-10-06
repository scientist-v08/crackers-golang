package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/scientist-v08/crackers/constants"
	"github.com/scientist-v08/crackers/controller"
	"github.com/scientist-v08/crackers/middleware"
)

func RegisterInventoryRoutes(r *fiber.App) {
	r.Post("/post/inventory", middleware.RequireAnyRole(constants.RoleAdmin), controller.AddItemToInventoryHandler)
	r.Get("/get/inventory", middleware.RequireAnyRole(constants.RoleAdmin), controller.GetPaginatedInventoryItems)
	r.Post("/post/updateInventory", middleware.RequireAnyRole(constants.RoleAdmin), controller.UpdateInventory)
}