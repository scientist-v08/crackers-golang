package service

import (
	"errors"
	"os"
	"slices"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/scientist-v08/crackers/constants"
	"github.com/scientist-v08/crackers/model"
	"github.com/scientist-v08/crackers/repository"
	"github.com/scientist-v08/crackers/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignUp(email, password string) error {
	// Check if user already exists
	existing, err := repository.FindByEmail(email)
	if err == nil && existing != nil {
		return errors.New("user already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if !utils.IsPasswordValid(password) {
		return errors.New("password must contain at least 1 uppercase letter, 1 lowercase letter, 1 special character, and be at least 8 characters long")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return errors.New("failed to hash password")
	}

	newUser := &model.User{
		Email:    email,
		Password: string(hash),
		Roles:    []string{"ROLE_USER"},
	}

	return repository.CreateUser(newUser)
}

func AdminSignUp(email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return errors.New("failed to hash password")
	}

	newUser := &model.User{
		Email:    email,
		Password: string(hash),
		Roles:    []string{"ROLE_ADMIN", "ROLE_USER"},
	}

	return repository.CreateUser(newUser)
}

func Login(email, password string) (string, []model.Routes, error) {
	user, err := repository.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("invalid email ID")
		}
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Roles,
		"exp": time.Now().Add(time.Hour * 12).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", nil, errors.New("failed to create JWT token")
	}

	isAdmin := slices.Contains(user.Roles, "ROLE_ADMIN")
	role := constants.RoleUser
	if isAdmin {
		role = constants.RoleAdmin
	}

	routes, err := repository.FindByRole(role)
	if err != nil {
		return "", nil, err
	}

	return tokenString, routes, nil
}