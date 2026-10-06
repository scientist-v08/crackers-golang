package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/scientist-v08/crackers/dto"
	"github.com/scientist-v08/crackers/service"
)

// ---------- Gin Handler -------------
func CreateExpense(c fiber.Ctx) error {

	// 1. Obtain the request body
	var req dto.ExpenseReqBody
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Incorrect request body",
		})
	}

	// 2. Save the data in the DB
	if err := service.CreateExpenseService(req); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// 3. Send success message
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"Success": "Item saved to inventory"})
}

func GetAllExpenses(c fiber.Ctx) error {
	// 1. Create variables
	var res dto.GetExpenseRes
	var err error

	// 2. Obtain results
	res, err = service.GetAllExpenseService()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	// 3. Return results
	return c.Status(fiber.StatusOK).JSON(res)
}
