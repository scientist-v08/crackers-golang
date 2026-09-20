package repository

import (
	"github.com/scientist-v08/crackers/initializers"
	"github.com/scientist-v08/crackers/model"
)

func FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := initializers.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *model.User) error {
	return initializers.DB.Create(user).Error
}