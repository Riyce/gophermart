package accrual

import (
	"database/sql"
	"github.com/riyce/gophermart/internal/models"
)

type DBOrder interface {
	UpdateOrderStatus(orderID string) error
	ProcessOrder(order order) error
	GetUnprocessedOrders() ([]order, error)
}

type DBOrderUpdater struct {
	db *sql.DB
}

func NewDBOrderUpdater(db *sql.DB) *DBOrderUpdater {
	return &DBOrderUpdater{db: db}
}

func (o *DBOrderUpdater) GetUnprocessedOrders() ([]order, error) {
	var ordersList []order
	rows, err := o.db.Query(getUnprocessedOrdersQuery, models.New, models.Processing)
	if err != nil || rows.Err() != nil {
		if err != nil {
			return ordersList, err
		} else {
			return ordersList, rows.Err()
		}
	}

	for rows.Next() {
		var order order
		if err := rows.Scan(&order.Order, &order.Status); err != nil {
			return ordersList, err
		}

		ordersList = append(ordersList, order)
	}
	return ordersList, nil
}

func (o *DBOrderUpdater) UpdateOrderStatus(orderID string) error {
	row, err := o.db.Query(setProcessingStatusQuery, models.Processing, orderID)
	if err != nil || row.Err() != nil {
		if err != nil {
			return err
		} else {
			return row.Err()
		}
	}

	return nil
}

func (o *DBOrderUpdater) ProcessOrder(order order) error {
	var userID int
	row := o.db.QueryRow(setProcessedStatusQuery, order.Status, order.Accrual, order.Order)
	if row.Err() != nil {
		return row.Err()
	}

	if err := row.Scan(&userID); err != nil {
		return err
	}
	if order.Accrual > 0 {
		_, err := o.db.Exec(addPointsToUser, order.Accrual, userID)
		return err
	}
	return nil
}
