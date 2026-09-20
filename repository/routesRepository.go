package repository

import (
	"github.com/scientist-v08/crackers/initializers"
	"github.com/scientist-v08/crackers/model"
)

func FindByRole(role string) ([]model.Routes, error) {
	var routes []model.Routes
	err := initializers.DB.Where("role = ?", role).Find(&routes).Error
	return routes, err
}