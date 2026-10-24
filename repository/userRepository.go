package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/scientist-v08/crackers/db"
	"github.com/scientist-v08/crackers/initializers"
)

func FindByEmail(email string) (*db.User, error) {
	user, err := initializers.Queries.FindUserByEmail(context.Background(), pgtype.Text{
		String: email,
		Valid:  true,
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(email, password string, roles []string) error {
	rolesJSON, errRs := json.Marshal(roles)
	if errRs != nil {
		return errRs
	}
	_, err := initializers.Queries.CreateUser(context.Background(), db.CreateUserParams{
		Email: pgtype.Text{
			String: email,
			Valid:  true,
		},
		Password: pgtype.Text{
			String: password,
			Valid:  true,
		},
		Roles: pgtype.Text{
			String: string(rolesJSON),
			Valid:  true,
		},
	})
	return err
}