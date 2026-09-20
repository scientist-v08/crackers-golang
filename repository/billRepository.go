package repository

import (
	"github.com/scientist-v08/crackers/initializers"
	"github.com/scientist-v08/crackers/model"
	"github.com/scientist-v08/crackers/utils"
	"gorm.io/gorm"
)

func AddNewBill(tx *gorm.DB, req *model.BillDetails) error {
	// Check for transaction
	db := initializers.DB
	if tx != nil {
		db = tx
	}

	// Now proceed with adding to DB
	// Build the Billing record + nested Purchases
	purchases := make([]model.Purchases, 0, len(req.BillItems))
	for _, it := range req.BillItems {
		purchases = append(purchases, model.Purchases{
			Mobile:   req.Mobile,
			MrpOrNet: utils.ToCents(it.MRPOrNet),
			Item:     it.Item,
			Quantity: int32(it.Quantity),
			Discount: it.Discount,
			SubTotal: utils.ToCents(it.SubTotal),
			// BillingID is set automatically by GORM
		})
	}
	// Decide FinalizedAmt
	finalizedAmt := req.FinalizedAmt
	if finalizedAmt == 0 {
		finalizedAmt = utils.ToCents(req.GrandTotal)
	}
	billing := model.Billing{
		User:         req.User,
		Mobile:       req.Mobile,
		GrandTotal:   utils.ToCents(req.GrandTotal),
		FinalizedAmt: finalizedAmt,
		Purchases:    purchases, // nested association
	}
	if errBilling := db.Create(&billing).Error; errBilling != nil {
		return errBilling
	}

	return nil
}