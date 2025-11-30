package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/scientist-v08/crackers/constants"
	"github.com/scientist-v08/crackers/controller"
	"github.com/scientist-v08/crackers/initializers"
	"github.com/scientist-v08/crackers/middleware"
)

func RegisterBillingRoutes(r *gin.Engine) {
	r.POST("/billing", middleware.RequireAnyRole(constants.RoleUser, constants.RoleAdmin), controller.CreateBillHandler(initializers.DB))
	r.POST("/billing/preview", middleware.RequireAnyRole(constants.RoleUser, constants.RoleAdmin), controller.CreatePreviewBillHandler)
}
