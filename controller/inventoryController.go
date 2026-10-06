package controller

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/scientist-v08/crackers/dto"
	"github.com/scientist-v08/crackers/service"
)

// ---------- Gin Handler -------------
func AddItemToInventoryHandler(c fiber.Ctx) error {

	// 1. Obtain the request body
	var reqBody dto.InventoryReqBody
	if err := c.Bind().JSON(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to obtain the request body"})
	}

	// 2. Validate the field State
	if err := reqBody.State.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	// 3. Save the data in the DB
	_, isError := service.AddItemToinventory(reqBody)
	if isError != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": isError.Error()})
	}

	// 4. Send success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"Success": "Item saved to inventory"})
}

func GetPaginatedInventoryItems(c fiber.Ctx) error {
	// Parse query params only
	pageNumber, _ := strconv.Atoi(c.Query("pageNumber", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	req := dto.InventoryListRequest{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		SortBy:     "id",
		Order:      "DESC",
		Title:      strings.TrimSpace(c.Query("title")),
		State:      strings.TrimSpace(c.Query("state")),
	}

	resp, err := service.GetPaginatedInventoryItems(req)

	if err != nil {
		// Map domain errors to HTTP status
		if strings.Contains(err.Error(), "invalid state") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func UpdateInventory(c fiber.Ctx) error {
	var req dto.UpdateInventoryRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if err := service.UpdateInventory(req); err != nil {
		switch {
		case errors.Is(err, service.ErrInventoryNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, service.ErrInvalidStateTransition):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update inventory"})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Inventory updated successfully",
	})
}
