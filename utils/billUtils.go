package utils

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/scientist-v08/crackers/model"
)

func ToCents(v float64) int32 {
	return int32(v)
}

func ParseDiscount(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty discount")
	}

	// Remove trailing %
	s = strings.TrimSuffix(s, "%")
	s = strings.TrimSpace(s)

	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	// If the original value was ≥ 1 we treat it as a percentage
	// (75 → 0.25, 10 → 0.90, etc.)
	if val >= 1 {
		return 1 - (val / 100), nil
	}

	// Already a multiplier (0.25, 0.9, …)
	return val, nil
}

// ---------- Helper: Re-calculate the sub-totals and the grand total ------------
func RecalculateSubTotalsAndGrandTotal(bill *model.BillDetails) error {
	var grandTotal float64
	discountMultiplier, err := ParseDiscount(bill.Discount)
	if err != nil {
		return fmt.Errorf("invalid discount %q: %w", bill.Discount, err)
	}

	for i := range bill.BillItems {
		item := &bill.BillItems[i]

		qty := float64(item.Quantity)
		mrp := item.MRPOrNet

		var subTotal float64

		if strings.HasPrefix(item.Item, "Other") {
			// Rule 1: Other* → MRP × Quantity
			subTotal = mrp * qty
		} else {
			// Rule 2: everything else → MRP × Quantity × Discount
			subTotal = mrp * qty * discountMultiplier
		}

		// Overwrite the value the client sent
		item.SubTotal = subTotal
		grandTotal += subTotal
	}

	// Overwrite the value the client sent
	bill.GrandTotal = grandTotal
	return nil
}

// ---------- Generate PDF ----------
func GeneratePDF(bill model.BillDetails, billID uint64) ([]byte, error) {
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
	f.SetFillColor(217, 221, 220)
	f.SetFont("Arial", "B", 10)
	headers := []string{"Sl.No", "Item", "MRP/Net", "Quantity", "SubTotal w/o Discount", "Discount", "SubTotal"}
	colWidths := []float64{15, 58, 23, 18, 28, 23, 25}
	headerHeight := 12.0
	for i, h := range headers {
		// Draw fixed-height bordered cell first, then overlay text → uniform height for ALL headers
		x := f.GetX()
		y := f.GetY()

		// 1. Draw empty tall cell with border (forces uniform height)
		f.CellFormat(colWidths[i], headerHeight, "", "1", 0, "", true, 0, "")

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
		f.CellFormat(18, 8, strconv.Itoa(int(item.Quantity)), "1", 0, "C", false, 0, "")
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
		f.SetFont("Arial", "B", 12) // slightly bigger for emphasis
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