package service

import (
	"context"
	"fmt"

	"github.com/scientist-v08/crackers/db"
	"github.com/scientist-v08/crackers/dto"
	"github.com/scientist-v08/crackers/repository"
)

func SaveTempBill(ctx context.Context, req dto.TempBillDetails) error {
	tx, err := repository.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	finalizedAmt := req.GrandTotal
	if req.FinalizedAmt > 0 {
		finalizedAmt = req.FinalizedAmt
	}
	billingID, err := repository.InsertTempBilling(ctx, tx, db.InsertTempBillingParams{
		User:         req.User,
		Mobile:       req.Mobile,
		Comments:     req.Comments,
		GrandTotal:   req.GrandTotal,
		FinalizedAmt: finalizedAmt,
	})
	if err != nil {
		return fmt.Errorf("insert billing: %w", err)
	}

	purchaseParams := make([]db.InsertTempPurchasesParams, 0, len(req.TempBillItems))
	for _, item := range req.TempBillItems {
		purchaseParams = append(purchaseParams, db.InsertTempPurchasesParams{
			BillingID: billingID,
			Mobile:    req.Mobile,
			MrpOrNet:  item.MrpOrNet,
			Item:      item.Item,
			Quantity:  item.Quantity,
			Discount:  item.Discount,
			SubTotal:  item.SubTotal,
		})
	}

	if _, err := repository.InsertTempPurchases(ctx, tx, purchaseParams); err != nil {
		return fmt.Errorf("insert purchases: %w", err)
	}

	return tx.Commit(ctx)
}

func GetAndConsumeTempBill(ctx context.Context, id int16) (*dto.TempBillResponse, error) {
	tx, err := repository.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Sequential is safer with a single transaction
	billing, err := repository.GetTempBillingByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("get billing: %w", err)
	}

	purchases, err := repository.ListTempPurchasesByBillingID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("get purchases: %w", err)
	}

	// Cascade will delete related purchases
	if err := repository.DeleteTempBillingByID(ctx, tx, id); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	items := make([]dto.TempBillItem, 0, len(purchases))
	for _, p := range purchases {
		items = append(items, dto.TempBillItem{
			MrpOrNet: p.MrpOrNet,
			Item:     p.Item,
			Quantity: p.Quantity,
			Discount: p.Discount,
			SubTotal: p.SubTotal,
		})
	}

	return &dto.TempBillResponse{
		ID:           billing.ID,
		User:         billing.User,
		Mobile:       billing.Mobile,
		Comments:     billing.Comments,
		GrandTotal:   billing.GrandTotal,
		FinalizedAmt: billing.FinalizedAmt,
		Items:        items,
	}, nil
}

func GetTempBillingCommentsAndFinalized(ctx context.Context) ([]dto.CommentsAndFinalizedAmt, error) {
	rows, err := repository.GetTempBillingCommentsAndFinalizedByID(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.CommentsAndFinalizedAmt, 0, len(rows))
	for _, row := range rows {
		result = append(result, dto.CommentsAndFinalizedAmt{
			ID:           row.ID,
			Comments:     row.Comments,
			FinalizedAmt: row.FinalizedAmt,
		})
	}
	return result, nil
}

func DeleteTempBill(ctx context.Context, id int16) error {
	tx, err := repository.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Cascade will automatically delete related purchases
	if err := repository.DeleteTempBillingByID(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit(ctx)
}