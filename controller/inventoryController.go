package controller

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/scientist-v08/crackers/dto"
	"github.com/scientist-v08/crackers/model"
	"github.com/scientist-v08/crackers/service"
)

// ---------- Gin Handler -------------
func AddItemToInventoryHandler(c *gin.Context) {

	// 1. Obtain the request body
	var reqBody model.InventoryReqBody
	if err := c.ShouldBindBodyWithJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to obtain the request body"})
		return
	}

	// 2. Validate the field State
	if err := reqBody.State.Validate(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// 3. Save the data in the DB
	_, isError := service.AddItemToinventory(reqBody)
	if isError != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": isError.Error()})
		return
	}

	// 4. Send success message
	c.JSON(http.StatusOK, gin.H{"Success": "Item saved to inventory"})
}

func GetPaginatedInventoryItems(c *gin.Context) {
	// Parse query params only
	pageNumber, _ := strconv.Atoi(c.DefaultQuery("pageNumber", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	req := dto.InventoryListRequest{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		SortBy:     c.DefaultQuery("sortBy", "id"),
		Order:      strings.ToUpper(c.DefaultQuery("order", "DESC")),
		Title:      strings.TrimSpace(c.Query("title")),
		State:      strings.TrimSpace(c.Query("state")),
	}

	resp, err := service.GetPaginatedInventoryItems(req)

	if err != nil {
		// Map domain errors to HTTP status
		if strings.Contains(err.Error(), "invalid state") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func UpdateInventory(c *gin.Context) {
	var req model.UpdateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := service.UpdateInventory(req); err != nil {
		switch {
		case errors.Is(err, service.ErrInventoryNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrInvalidStateTransition):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Inventory updated successfully",
	})
}
