package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/scientist-v08/crackers/db"
	"github.com/scientist-v08/crackers/dto"
	"github.com/scientist-v08/crackers/initializers"
	"github.com/scientist-v08/crackers/utils"
)

func AddNewBill(ctx context.Context, tx pgx.Tx, req *dto.BillDetails) (int64, error) {
	// Check for transaction
	q := initializers.Queries
	if tx != nil {
		q = initializers.Queries.WithTx(tx)
	}

	// 1. Decide FinalizedAmt
	finalizedAmt := req.FinalizedAmt
	if finalizedAmt == 0 {
		finalizedAmt = utils.ToCents(req.GrandTotal)
	}

	// 2. Create the Billing record
	billing, err := q.CreateBilling(ctx, db.CreateBillingParams{
		User: pgtype.Text{
			String: req.User,
			Valid:  true,
		},
		Mobile: pgtype.Text{
			String: req.Mobile,
			Valid:  true,
		},
		GrandTotal: pgtype.Int4{
			Int32: utils.ToCents(req.GrandTotal),
			Valid: true,
		},
		FinalizedAmt: pgtype.Int4{
			Int32: finalizedAmt,
			Valid: true,
		},
	})
	if err != nil {
		return 0, err
	}

	// 3. Prepare bulk purchases
	purchases := make([]db.CreatePurchasesParams, 0, len(req.BillItems))
	for _, it := range req.BillItems {
		purchases = append(purchases, db.CreatePurchasesParams{
			BillingID: billing.ID,
			Mobile:    pgtype.Text{String: req.Mobile, Valid: true},
			MrpOrNet:  pgtype.Int4{Int32: utils.ToCents(it.MRPOrNet), Valid: true},
			Item:      pgtype.Text{String: it.Item, Valid: true},
			Quantity:  pgtype.Int4{Int32: int32(it.Quantity), Valid: true},
			Discount:  pgtype.Text{String: it.Discount, Valid: true},
			SubTotal:  pgtype.Int4{Int32: utils.ToCents(it.SubTotal), Valid: true},
		})
	}

	// 4. Bulk insert all purchases at once
	if len(purchases) > 0 {
		_, err = q.CreatePurchases(ctx, purchases)
		if err != nil {
			return 0, err
		}
	}

	return billing.ID, nil
}