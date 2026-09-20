package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scientist-v08/crackers/dto"
	"github.com/scientist-v08/crackers/initializers"
	"github.com/scientist-v08/crackers/model"
)

func GetAllPrices(c *gin.Context) {
	var allPrices []model.PriceList
	if err := initializers.DB.Find(&allPrices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	grouped := make(dto.PriceListBrandMapping)
	for _, item := range allPrices {
		grouped[item.Brand] = append(grouped[item.Brand], dto.Products{
			Id: item.Id,
			Price: item.Price,
			Item: item.Item,
		})
	}

	result := make([]dto.ProductsList, 0, len(grouped))
	for brand, product := range grouped {
		result = append(result, dto.ProductsList{
			Brand: brand,
			List: product,
		})
	}

	c.JSON(http.StatusOK, result)
}