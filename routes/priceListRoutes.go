package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/scientist-v08/crackers/constants"
	"github.com/scientist-v08/crackers/controller"
	"github.com/scientist-v08/crackers/middleware"
)

func RegisterPriceList(r *gin.Engine) {
	r.GET("/get/price-list", middleware.RequireAnyRole(constants.RoleUser, constants.RoleAdmin), controller.GetAllPrices)
}