package db

import (
	"database/sql"
	"github.com/riyce/gophermart/internal/models"
)

type AuthDB interface {
	CreateUser(user *models.User) error
	GetUser(user *models.User) error
	GetUserID(apiKey string) (int, error)
}

type OrderDB interface {
	AddOrder(orderID string, userID int) error
	GetOrders(userID int) ([]models.Order, error)
}

type BalanceDB interface {
	GetBalance(userID int) (models.Balance, error)
	Withdraw(userID int, withdraw models.Withdraw) error
	GetWithdrawsHistory(userID int) ([]models.Withdraw, error)
}

type DB struct {
	AuthDB
	BalanceDB
	OrderDB
}

func NewDB(db *sql.DB) *DB {
	return &DB{
		AuthDB:    NewAuthController(db),
		OrderDB:   NewOrdersController(db),
		BalanceDB: NewBalanceController(db),
	}
}
