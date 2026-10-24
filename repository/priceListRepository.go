package repository

import (
	"context"

	"github.com/scientist-v08/crackers/db"
	"github.com/scientist-v08/crackers/initializers"
)

func GetAllPrices() ([]db.PriceList, error) {
	allPrices, err := initializers.Queries.GetAllPriceList(context.Background())
	if err != nil {
		return nil, err
	}
	return allPrices, nil
}