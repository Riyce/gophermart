package services

import (
	"github.com/riyce/gophermart/internal/db"
	"github.com/riyce/gophermart/internal/models"
)

type Auth interface {
	CreateUser(user *models.User) error
	GetUser(user *models.User) error
	GetUserID(apiKey string) (int, error)
}

type Order interface {
	AddOrder(orderID, userID int) error
	GetOrders(userID int) ([]models.Order, error)
}

type Balance interface {
	GetBalance(userID int) (models.Balance, error)
	Withdraw(userID int, withdraw models.Withdraw) error
	GetWithdrawsHistory(userID int) ([]models.Withdraw, error)
}

type Service struct {
	Auth
	Balance
	Order
}

func NewService(db *db.DB, key string) *Service {
	return &Service{
		Auth:    NewAuthService(db.AuthDB, key),
		Order:   NewOrderService(db.OrderDB),
		Balance: NewBalanceService(db.BalanceDB),
	}
}
