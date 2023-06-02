package db

import (
	"database/sql"
	"fmt"
	"github.com/riyce/gophermart/internal/models"
	"github.com/riyce/gophermart/internal/utils"
	"github.com/rs/zerolog/log"
)

const (
	ordersTableName string = "orders"
	ordersDB        string = "Orders DB"
)

type OrdersController struct {
	db *sql.DB
}

func NewOrdersController(db *sql.DB) *OrdersController {
	return &OrdersController{db: db}
}

func (o *OrdersController) AddOrder(orderID string, userID int) error {
	query := fmt.Sprintf(addOrderQuery, ordersTableName)
	row, createErr := o.db.Query(query, orderID, userID, models.New)
	if createErr != nil || row.Err() != nil {
		var ownerID int
		secondQuery := fmt.Sprintf(getUserIDByOrderNumberQuery, ordersTableName)
		row := o.db.QueryRow(secondQuery, orderID)
		if err := row.Scan(&ownerID); err != nil {
			log.Error().
				Err(err).
				Str("service", ordersDB).
				Msg(fmt.Sprintf("error scan order's owner ID for order %s", orderID))
			return utils.ErrSomethingWentWrong
		}

		if userID == ownerID {
			log.Warn().
				Str("service", ordersDB).
				Msg(fmt.Sprintf("order %s already exists", orderID))
			return utils.ErrUsersOrderAlreadyExists
		} else {
			log.Warn().
				Str("service", ordersDB).
				Msg(fmt.Sprintf("order %s already exists from user %d", orderID, ownerID))
			return utils.ErrOtherOrderAlreadyExists
		}
	}

	return nil
}

func (o *OrdersController) GetOrders(userID int) ([]models.Order, error) {
	var ordersList []models.Order
	query := fmt.Sprintf(getUsersOrdersQuery, ordersTableName)
	rows, err := o.db.Query(query, userID)
	if err != nil || rows.Err() != nil {
		if err != nil {
			log.Error().
				Err(err).
				Str("service", ordersDB).
				Msg(fmt.Sprintf("error on get orders for user %d", userID))
		} else {
			log.Error().
				Err(rows.Err()).
				Str("service", ordersDB).
				Msg(fmt.Sprintf("error on get orders for user %d", userID))
		}
		return ordersList, utils.ErrSomethingWentWrong
	}

	defer rows.Close()

	for rows.Next() {
		var order models.Order
		var accrual sql.NullFloat64
		err := rows.Scan(&order.Number, &order.Status, &accrual, &order.CreatedAt)
		if accrual.Valid {
			order.Accrual = float32(accrual.Float64)
		}
		if err != nil {
			log.Error().
				Err(err).
				Str("service", ordersDB).
				Msg(fmt.Sprintf("error on scan orders for user %d", userID))
			return ordersList, utils.ErrSomethingWentWrong
		}
		ordersList = append(ordersList, order)
	}

	return ordersList, nil
}
