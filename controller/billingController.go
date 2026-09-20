package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/scientist-v08/crackers/model"
	"github.com/scientist-v08/crackers/service"
	"github.com/scientist-v08/crackers/utils"
)

// ---------- Gin handler ----------
func CreateBillHandler(c *gin.Context) {
	var req model.BillDetails
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Discard client totals and recalculate
	if err := utils.RecalculateSubTotalsAndGrandTotal(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save data into DB
	if err := service.GenerateBillService(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate PDF
	pdfBytes, err := service.PreviewBillService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send PDF in response
	loc, _ := time.LoadLocation("Asia/Kolkata")
	nowIST := time.Now().In(loc)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="bill_preview_%s.pdf"`, nowIST))
	c.Data(http.StatusCreated, "application/pdf", pdfBytes)
}

func CreatePreviewBillHandler(c *gin.Context) {
	// 1. Obtain request body
	var req model.BillDetails
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Generate PDF
	req.User = "NA--PREVIEW"
	req.Mobile = "NA--PREVIEW"

	// Discard client totals and recalculate
	if err := utils.RecalculateSubTotalsAndGrandTotal(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pdfBytes, err := service.PreviewBillService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send PDF in response
	loc, _ := time.LoadLocation("Asia/Kolkata")
	nowIST := time.Now().In(loc)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="bill_preview_%s.pdf"`, nowIST))
	c.Data(http.StatusCreated, "application/pdf", pdfBytes)
}
