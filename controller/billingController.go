package controller

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/scientist-v08/crackers/dto"
	"github.com/scientist-v08/crackers/service"
	"github.com/scientist-v08/crackers/utils"
)

// ---------- Gin handler ----------
func CreateBillHandler(c fiber.Ctx) error {
	var req dto.BillDetails
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Discard client totals and recalculate
	if err := utils.RecalculateSubTotalsAndGrandTotal(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Save data into DB
	_id := int64(0)
	var generateBillErr error
	if _id, generateBillErr = service.GenerateBillService(req); generateBillErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": generateBillErr.Error()})
	}

	// Generate PDF
	pdfBytes, err := service.PreviewBillService(req, _id)
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

func CreatePreviewBillHandler(c fiber.Ctx) error {
	// 1. Obtain request body
	var req dto.BillDetails
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// 2. Generate PDF
	req.User = "NA--PREVIEW"
	req.Mobile = "NA--PREVIEW"

	// 3. Discard client totals and recalculate
	if err := utils.RecalculateSubTotalsAndGrandTotal(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		
	}

	pdfBytes, err := service.PreviewBillService(req, 0)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// 4. Send PDF in response
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
