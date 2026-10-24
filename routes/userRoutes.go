package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/scientist-v08/crackers/controller"
)

func RegisterUserRoutes(r *fiber.App) {
	r.Post("/create/user", controller.SignUp)
	r.Post("/create/admin", controller.AdminSignUp)
	r.Post("/login/user", controller.Login)
}