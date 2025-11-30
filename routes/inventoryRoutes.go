package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/scientist-v08/crackers/constants"
	"github.com/scientist-v08/crackers/controller"
	"github.com/scientist-v08/crackers/initializers"
	"github.com/scientist-v08/crackers/middleware"
)

func RegisterInventoryRoutes(r *gin.Engine) {
	r.POST("/post/inventory", middleware.RequireAnyRole(constants.RoleAdmin), controller.AddItemToInventoryHandler)
	r.GET("/get/inventory", middleware.RequireAnyRole(constants.RoleAdmin), controller.GetPaginatedInventoryItems)
	r.POST("/post/updateInventory", middleware.RequireAnyRole(constants.RoleAdmin), controller.UpdateInventory(initializers.DB))
}