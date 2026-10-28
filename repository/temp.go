package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/scientist-v08/crackers/db"
	"github.com/scientist-v08/crackers/initializers"
)

func InsertTempBilling(ctx context.Context, tx pgx.Tx, arg db.InsertTempBillingParams) (int16, error) {
	q := db.New(tx)
	return q.InsertTempBilling(ctx, arg)
}

func InsertTempPurchases(ctx context.Context, tx pgx.Tx, args []db.InsertTempPurchasesParams) (int64, error) {
	q := db.New(tx)
	return q.InsertTempPurchases(ctx, args)
}

func GetTempBillingByID(ctx context.Context, tx pgx.Tx, id int16) (db.TempBilling, error) {
	q := db.New(tx)
	return q.GetTempBillingByID(ctx, id)
}

func ListTempPurchasesByBillingID(ctx context.Context, tx pgx.Tx, billingID int16) ([]db.TempPurchase, error) {
	q := db.New(tx)
	return q.ListTempPurchasesByBillingID(ctx, billingID)
}

func DeleteTempBillingByID(ctx context.Context, tx pgx.Tx, id int16) error {
	q := db.New(tx)
	return q.DeleteTempBillingByID(ctx, id)
}

func DeleteTempPurchasesByBillingID(ctx context.Context, tx pgx.Tx, billingID int16) error {
	q := db.New(tx)
	return q.DeleteTempPurchasesByBillingID(ctx, billingID)
}

func GetTempBillingCommentsAndFinalizedByID(ctx context.Context) ([]db.GetTempBillingCommentsAndFinalizedByIDRow, error) {
	q := initializers.Queries
	return q.GetTempBillingCommentsAndFinalizedByID(ctx)
}