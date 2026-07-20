package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/scientist-v08/crackers/model"
	"gorm.io/gorm"
)

// ---------- Front-end interfaces (mirrored in Go) ----------
type Item struct {
	SlNo     int     `json:"slNo"`
	Item     string  `json:"item" binding:"required"`
	MRPOrNet float64 `json:"mrpOrNet" binding:"required,gt=0"`
	Quantity int     `json:"quantity" binding:"required,gt=0"`
	Discount string  `json:"discount"`
	SubTotal float64 `json:"subTotal" binding:"required,gt=0"`
}

type BillDetails struct {
	User       	 string  `json:"user" binding:"required"`
	Mobile     	 string  `json:"mobile" binding:"required"`
	GrandTotal 	 float64 `json:"grandTotal" binding:"required,gt=0"`
	FinalizedAmt int32 `json:"finalizedAmt"`
	BillItems    []Item  `json:"billItems" binding:"required,min=1"`
}

// ---------- Helper: convert float → int32 (cents) ----------
func toCents(v float64) int32 {
	return int32(v)
}

// ---------- Generate PDF ----------
func generatePDF(bill BillDetails, billID uint) ([]byte, error) {
	f := gofpdf.New("P", "mm", "A4", "")
	f.AddPage()
	f.SetFont("Arial", "B", 16)

	// Heading
	title := "Vinayaka Traders"
	f.CellFormat(190, 10, title, "", 1, "C", false, 0, "")
	f.Ln(5)

	// Flex-like section
	y := f.GetY()

	// Left div: Customer details
	f.SetFont("Arial", "", 10)
	f.SetXY(10, y)
	customerText := fmt.Sprintf("Customer: %s\nMobile: %s\nBill ID: %d", bill.User, bill.Mobile, billID)
	f.MultiCell(90, 5, customerText, "", "L", false)
	yLeft := f.GetY()

	// Right div: Terms
	f.SetXY(110, y)
	termsText := "T&C: Quality not guaranteed by retailer. Contact brand for complaints."
	f.MultiCell(90, 6, termsText, "", "R", false)
	yRight := f.GetY()

	// Move to the bottom of the taller section
	maxY := yLeft
	if yRight > maxY {
		maxY = yRight
	}
	f.SetY(maxY)

	f.Ln(5)

	// Table headers
	f.SetFont("Arial", "B", 10)
	headers := []string{"Sl.No", "Item", "MRP/Net", "Quantity", "SubTotal w/o Discount", "Discount", "SubTotal"}
	colWidths := []float64{15, 58, 23, 18, 28, 23, 25}
	headerHeight := 12.0
	for i, h := range headers {
		// Draw fixed-height bordered cell first, then overlay text → uniform height for ALL headers
		x := f.GetX()
		y := f.GetY()
		
		// 1. Draw empty tall cell with border (forces uniform height)
		f.CellFormat(colWidths[i], headerHeight, "", "1", 0, "", false, 0, "")
		
		// 2. Reset position and write the text (supports \n)
		f.SetXY(x, y)
		f.MultiCell(colWidths[i], 6, h, "", "C", false)
		
		// Move to next column at original Y
		f.SetXY(x+colWidths[i], y)
	}
	f.Ln(headerHeight)

	// Table rows
	f.SetFont("Arial", "", 10)
	for _, item := range bill.BillItems {
		f.CellFormat(15, 8, strconv.Itoa(item.SlNo), "1", 0, "C", false, 0, "")
		f.CellFormat(58, 8, item.Item, "1", 0, "L", false, 0, "")
		f.CellFormat(23, 8, fmt.Sprintf("%.2f", item.MRPOrNet), "1", 0, "R", false, 0, "")
		f.CellFormat(18, 8, strconv.Itoa(item.Quantity), "1", 0, "C", false, 0, "")
		// New column: SubTotal without discount = MRPOrNet * Quantity
		subTotalNoDisc := item.MRPOrNet * float64(item.Quantity)
		f.CellFormat(28, 8, fmt.Sprintf("%.2f", subTotalNoDisc), "1", 0, "R", false, 0, "")
		f.CellFormat(23, 8, item.Discount, "1", 0, "C", false, 0, "")
		f.CellFormat(25, 8, fmt.Sprintf("%.2f", item.SubTotal), "1", 0, "R", false, 0, "")
		f.Ln(-1)
	}

	// Grand total
	f.Ln(8)
	f.SetFont("Arial", "B", 12)
	f.CellFormat(150, 10, "Grand Total:", "", 0, "R", false, 0, "")
	f.CellFormat(40, 10, fmt.Sprintf("%.2f", bill.GrandTotal), "", 0, "R", false, 0, "")
	f.Ln(4)

	// Finalized Amount
	if bill.FinalizedAmt > 0 {
		f.SetFont("Arial", "B", 12)  // slightly bigger for emphasis
		f.CellFormat(150, 12, "Finalized Amount:", "", 0, "R", false, 0, "")
		f.CellFormat(40, 12, fmt.Sprintf("%d", bill.FinalizedAmt), "", 0, "R", false, 0, "")
	}

	if err := f.Error(); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err := f.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ---------- Gin handler ----------
func CreateBillHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req BillDetails
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var createdBillingID uint
		// ---- Transaction ----
		errTrn := db.Transaction(func(tx *gorm.DB) error {

			// 1. Insert Billing row
			billing := model.Billing{
				User:       req.User,
				Mobile:     req.Mobile,
				GrandTotal: toCents(req.GrandTotal),
				FinalizedAmt: req.FinalizedAmt,
			}
			if errBilling := tx.Create(&billing).Error; errBilling != nil {
				return errBilling
			}
			createdBillingID = billing.ID
			
			// 2. Insert each item as a Purchases row
			var purchases []model.Purchases
			for _, it := range req.BillItems {
				purchases = append(purchases, model.Purchases{
					BillingID: billing.ID,
					Mobile:    req.Mobile,
					MrpOrNet:  toCents(it.MRPOrNet),
					Item:      it.Item,
					Quantity:  int32(it.Quantity),
					Discount:  it.Discount,
					SubTotal:  toCents(it.SubTotal),
				})
			}
			if errBulk := tx.Create(&purchases).Error; errBulk != nil {
				return errBulk
			}
			return nil
		})

		if errTrn != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errTrn.Error()})
			return
		}

		// Generate PDF
		pdfBytes, errPdf := generatePDF(req, createdBillingID)
		if errPdf != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errPdf.Error()})
			return
		}

		// Send PDF in response
		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="bill_%d.pdf"`, createdBillingID))
		c.Data(http.StatusCreated, "application/pdf", pdfBytes)
	}
}

func CreatePreviewBillHandler(c *gin.Context) {
	var req BillDetails
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.User = "NA--PREVIEW"
	req.Mobile = "NA--PREVIEW"
	// Generate PDF
	pdfBytes, errPdf := generatePDF(req, 0)
	if errPdf != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errPdf.Error()})
		return
	}

	// Send PDF in response
	loc, _ := time.LoadLocation("Asia/Kolkata")
	nowIST := time.Now().In(loc)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="bill_preview_%s.pdf"`, nowIST))
	c.Data(http.StatusCreated, "application/pdf", pdfBytes)
}
