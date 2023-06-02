package db

import (
	"database/sql"
	"fmt"
	"github.com/riyce/gophermart/internal/models"
	"github.com/riyce/gophermart/internal/utils"
	"github.com/rs/zerolog/log"
)

const (
	usersTableName string = "users"
	authDB         string = "Auth DB"
)

type AuthController struct {
	db *sql.DB
}

func NewAuthController(db *sql.DB) *AuthController {
	return &AuthController{db: db}
}

func (a *AuthController) CreateUser(user *models.User) error {
	var id int
	template := fmt.Sprintf(createUserQuery, usersTableName)
	row := a.db.QueryRow(template, user.Login, user.Password, user.APIKey)
	if row.Err() != nil {
		log.Error().
			Err(row.Err()).
			Str("service", authDB).
			Msg(fmt.Sprintf("error on create user %s with API-key %s", user.Login, user.APIKey))
		return utils.ErrUserAlreadyExists
	}

	if err := row.Scan(&id); err != nil {
		return utils.ErrSomethingWentWrong
	}

	return nil
}

func (a *AuthController) GetUser(user *models.User) error {
	template := fmt.Sprintf(getUserQuery, usersTableName)
	row := a.db.QueryRow(template, user.Login)
	if err := row.Scan(&user.Password, &user.APIKey); err != nil {
		log.Error().
			Err(row.Err()).
			Str("service", authDB).
			Msg(fmt.Sprintf("error on get user %s", user.Login))
		return utils.ErrWrongCredentials
	}

	return nil
}

func (a *AuthController) GetUserID(apiKey string) (int, error) {
	var id int
	template := fmt.Sprintf(getUserByKeyQuery, usersTableName)
	row := a.db.QueryRow(template, apiKey)
	if row.Err() != nil {
		log.Error().
			Err(row.Err()).
			Str("service", authDB).
			Msg(fmt.Sprintf("error on get user with API-key %s", apiKey))
	}

	if err := row.Scan(&id); err != nil {
		log.Error().
			Err(err).
			Str("service", authDB).
			Msg(fmt.Sprintf("error on get user with API-key %s", apiKey))
		return 0, utils.ErrWrongAPIKey
	}

	return id, nil
}
