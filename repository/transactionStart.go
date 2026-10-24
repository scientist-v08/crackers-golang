package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/scientist-v08/crackers/initializers"
)

func Begin(ctx context.Context) (pgx.Tx, error) {
	return initializers.Pool.Begin(ctx)
}