package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/scientist-v08/crackers/dto"
	"github.com/scientist-v08/crackers/service"
)

func SignUp(c fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to read request body"})
	}

	if err := service.SignUp(req.Email, req.Password); err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "user already exists" {
			status = fiber.StatusConflict
		}
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Creating a new user: Successful"})
}

func AdminSignUp(c fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		IsAdmin  bool   `json:"isAdmin"`
	}

	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to read request body"})
	}

	if !req.IsAdmin {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Only admins can use this API"})
	}

	if err := service.AdminSignUp(req.Email, req.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Creating a new user: Successful"})
}

func Login(c fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to read request body"})
	}

	token, routes, err := service.Login(req.Email, req.Password)
	if err != nil {
		status := fiber.StatusBadRequest
		if err.Error() == "invalid email ID" || err.Error() == "invalid password" {
			status = fiber.StatusUnauthorized // more accurate
		}
		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	}

	dtoRoutes := dto.ToRoutes(routes)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": token,
		"routes":       dtoRoutes,
	})
}