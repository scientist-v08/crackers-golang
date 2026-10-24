package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/scientist-v08/crackers/constants"
	"github.com/scientist-v08/crackers/controller"
	"github.com/scientist-v08/crackers/middleware"
)

func RegisterInventoryRoutes(r *fiber.App) {
	r.Post("/admin/inventory", middleware.RequireAnyRole(constants.RoleAdmin), controller.AddItemToInventoryHandler)
	r.Get("/admin/inventory", middleware.RequireAnyRole(constants.RoleAdmin), controller.GetPaginatedInventoryItems)
	r.Post("/admin/updateInventory", middleware.RequireAnyRole(constants.RoleAdmin), controller.UpdateInventory)
}