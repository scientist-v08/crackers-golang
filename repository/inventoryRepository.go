package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/scientist-v08/crackers/db"
	"github.com/scientist-v08/crackers/dto"
	"github.com/scientist-v08/crackers/initializers"
)

// AddItemToInventory inserts a new inventory record.
// Pass tx = nil if you don't need a transaction.
func AddItemToInventory(ctx context.Context, tx pgx.Tx, inventory dto.Inventory) (db.Inventory, error) {
	q := initializers.Queries
	if tx != nil {
		q = initializers.Queries.WithTx(tx)
	}

	return q.AddInventoryItem(ctx, db.AddInventoryItemParams{
		BrandOrCompany: pgtype.Text{String: inventory.BrandOrCompany, Valid: true},
		State:          pgtype.Text{String: string(inventory.State), Valid: true},
		Item:           pgtype.Text{String: inventory.Item, Valid: true},
		NumOfBoxes:     pgtype.Int4{Int32: inventory.NumOfBoxes, Valid: true},
		NumOfCartons:   pgtype.Int4{Int32: inventory.NumOfCartons, Valid: true},
		PricePerCarton: pgtype.Int4{Int32: inventory.PricePerCarton, Valid: true},
		SubTotal:       pgtype.Int4{Int32: inventory.SubTotal, Valid: true},
	})
}

// FindPaginatedInventory – sorting is hardcoded to id DESC (as you decided)
func FindPaginatedInventory(ctx context.Context, offset, limit int32, filters dto.InventoryFilter) ([]db.Inventory, error) {
	return initializers.Queries.FindPaginatedInventory(ctx, db.FindPaginatedInventoryParams{
		Column1: filters.State, // state filter (empty string = no filter)
		Column2: filters.Title, // title/item filter
		Limit:   limit,
		Offset:  offset,
	})
}

// CountInventory
func CountInventory(ctx context.Context, filters dto.InventoryFilter) (int64, error) {
	return initializers.Queries.CountInventory(ctx, db.CountInventoryParams{
		Column1: filters.State,
		Column2: filters.Title,
	})
}

// SumInventorySubTotal
func SumInventorySubTotal(ctx context.Context, filters dto.InventoryFilter) (int64, error) {
	return initializers.Queries.SumInventorySubTotal(ctx, db.SumInventorySubTotalParams{
		Column1: filters.State,
		Column2: filters.Title,
	})
}

// FindByInvId
func FindByInvId(ctx context.Context, tx pgx.Tx, id int64) (*db.Inventory, error) {
	q := initializers.Queries
	if tx != nil {
		q = initializers.Queries.WithTx(tx)
	}

	inv, err := q.FindInventoryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

// UpdateInventory – now uses the two dedicated queries
func UpdateInventory(ctx context.Context, tx pgx.Tx, id int64, updates any) error {
	q := initializers.Queries
	if tx != nil {
		q = initializers.Queries.WithTx(tx)
	}

	switch u := updates.(type) {

	case dto.UpdateInventoryState:
		return q.UpdateInventoryState(ctx, db.UpdateInventoryStateParams{
			ID: id,
			State: pgtype.Text{
				String: string(u.State),
				Valid:  true,
			},
			SubTotal: pgtype.Int4{
				Int32: u.SubTotal,
				Valid: true,
			},
		})

	case dto.InventoryUpdateDiff:
		return q.UpdateInventoryDiff(ctx, db.UpdateInventoryDiffParams{
			ID: id,
			NumOfCartons: pgtype.Int4{
				Int32: u.NumOfCartons,
				Valid: true,
			},
			SubTotal: pgtype.Int4{
				Int32: u.SubTotal,
				Valid: true,
			},
		})

	default:
		return errors.New("invalid update type received")
	}
}