package service

import (
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/scientist-v08/crackers/constants"
	"github.com/scientist-v08/crackers/db"
	"github.com/scientist-v08/crackers/initializers"
	"github.com/scientist-v08/crackers/repository"
	"github.com/scientist-v08/crackers/utils"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(email, password string) error {
	// Check if user already exists
	existing, err := repository.FindByEmail(email)
	if err == nil && existing != nil {
		return errors.New("user already exists")
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
        return err
    }

	if !utils.IsPasswordValid(password) {
		return errors.New("password must contain at least 1 uppercase letter, 1 lowercase letter, 1 special character, and be at least 8 characters long")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return errors.New("failed to hash password")
	}

	return repository.CreateUser(email, string(hash), []string{"ROLE_USER"})
}

func AdminSignUp(email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return errors.New("failed to hash password")
	}

	return repository.CreateUser(email, string(hash), []string{"ROLE_ADMIN", "ROLE_USER"})
}

func Login(email, password string) (string, []db.Route, error) {
	user, err := repository.FindByEmail(email)
	if err != nil {
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(password)); err != nil {
		return "", nil, errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Roles,
		"exp": time.Now().Add(time.Hour * 12).Unix(),
	})

	secret := initializers.Secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", nil, errors.New("failed to create JWT token")
	}

	var roles []string

	// Trim any surrounding whitespace first
	raw := strings.TrimSpace(user.Roles.String)
	if err := json.Unmarshal([]byte(raw), &roles); err != nil {
		return "", nil, err
	}

	isAdmin := slices.Contains(roles, "ROLE_ADMIN")
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