package repository

import (
	"context"

	"github.com/scientist-v08/crackers/db"
	"github.com/scientist-v08/crackers/initializers"
)

func FindByRole(role string) ([]db.Route, error) {
	return initializers.Queries.FindRoutesByRole(context.Background(), role)
}