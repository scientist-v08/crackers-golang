package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scientist-v08/crackers/initializers"
	"github.com/scientist-v08/crackers/model"
)

type Products struct {
	Id    uint   `json:"id"`
	Price uint   `json:"price"`
	Item  string `json:"item"`
}

type ProductsList struct {
	Brand string     `json:"brand"`
	List  []Products `json:"list"`
}

type PriceListBrandMapping map[string][]Products

func GetAllPrices(c *gin.Context) {
	var allPrices []model.PriceList
	if err := initializers.DB.Find(&allPrices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	grouped := make(PriceListBrandMapping)
	for _, item := range allPrices {
		grouped[item.Brand] = append(grouped[item.Brand], Products{
			Id: item.Id,
			Price: item.Price,
			Item: item.Item,
		})
	}

	result := make([]ProductsList, 0, len(grouped))
	for brand, product := range grouped {
		result = append(result, ProductsList{
			Brand: brand,
			List: product,
		})
	}

	c.JSON(http.StatusOK, result)
}