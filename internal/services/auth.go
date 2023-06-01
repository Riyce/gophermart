package services

import (
	"fmt"
	"github.com/riyce/gophermart/internal/db"
	"github.com/riyce/gophermart/internal/models"
	"github.com/riyce/gophermart/internal/utils"
	"github.com/rs/zerolog/log"
)

type AuthService struct {
	db  db.AuthDB
	key string
}

func NewAuthService(authDB db.AuthDB, key string) *AuthService {
	return &AuthService{db: authDB, key: key}
}

func (s *AuthService) CreateUser(user *models.User) error {
	passwordHash := utils.GetHash(user.Password, s.key)
	user.Password = passwordHash
	apiKey := utils.GenerateAPIKey()
	user.APIKey = apiKey
	if err := s.db.CreateUser(user); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) GetUser(user *models.User) error {
	password := user.Password
	if err := s.db.GetUser(user); err != nil {
		return err
	}

	equal := utils.CheckPasswordHash(password, user.Password, s.key)
	if !equal {
		log.Warn().
			Str("service", "Auth service").
			Msg(fmt.Sprintf("user %s used wrong credentials", user.Login))
		return utils.ErrWrongCredentials
	}

	return nil
}

func (s *AuthService) GetUserID(apiKey string) (int, error) {
	return s.db.GetUserID(apiKey)
}
