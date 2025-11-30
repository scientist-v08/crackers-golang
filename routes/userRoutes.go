package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/scientist-v08/crackers/controller"
)

func RegisterUserRoutes(r *gin.Engine) {
	r.POST("/create/user", controller.SignUp)
	r.POST("/create/admin", controller.AdminSignUp)
	r.POST("/login/user", controller.Login)
}