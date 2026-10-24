package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/scientist-v08/crackers/service"
)

func GetAllPrices(c fiber.Ctx) error {
	result, err := service.GetAllPricesService()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}