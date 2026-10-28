package controller

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/scientist-v08/crackers/dto"
	"github.com/scientist-v08/crackers/service"
	"github.com/scientist-v08/crackers/utils"
)

func SaveTempBill(c fiber.Ctx) error {
	var req dto.TempBillDetails
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	// Discard client totals and recalculate
	if err := utils.RecalculateSubTotalsAndGrandTotalForTemp(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Save data into DB
	if err := service.SaveTempBill(c, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	// Generate PDF
	pdfBytes, err := utils.GeneratePDFTemp(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Send PDF in response
	loc, _ := time.LoadLocation("Asia/Kolkata")
	nowIST := time.Now().In(loc)
	c.Set("Content-Disposition", fmt.Sprintf(
		`attachment; filename="bill_preview_%s.pdf"`,
		nowIST.Format("20060102_150405"),
	))

	return c.Status(fiber.StatusCreated).
		Type("pdf").
		Send(pdfBytes)
}

func GetSavedTempBill(c fiber.Ctx) error {
	var req dto.TempBillIDRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid request body",
		})
	}

	result, err := service.GetAndConsumeTempBill(c, req.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func DeleteSavedTempBill(c fiber.Ctx) error {
	var req dto.TempBillIDRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid request body",
		})
	}

	if err := service.DeleteTempBill(c, req.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func GetTempBillingCommentsAndFinalized(c fiber.Ctx) error {
	result, err := service.GetTempBillingCommentsAndFinalized(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}