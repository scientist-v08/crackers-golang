package service

import (
	"github.com/scientist-v08/crackers/dto"
	"github.com/scientist-v08/crackers/repository"
)

func GetAllPricesService() ([]dto.ProductsList, error) {
	allPrices, err := repository.GetAllPrices()
	if err != nil {
		return nil, err
	}

	grouped := make(dto.PriceListBrandMapping)

	for _, item := range allPrices {
		grouped[item.Brand.String] = append(grouped[item.Brand.String], dto.Products{
			Id:    uint(item.ID),
			Price: uint(item.Price.Int64),
			Item:  item.Item.String,
		})
	}

	result := make([]dto.ProductsList, 0, len(grouped))

	for brand, products := range grouped {
		result = append(result, dto.ProductsList{
			Brand: brand,
			List:  products,
		})
	}

	return result, nil
}