package services

import (
	"fmt"
	"github.com/riyce/gophermart/internal/db"
	"github.com/riyce/gophermart/internal/models"
	"github.com/riyce/gophermart/internal/utils"
	"github.com/rs/zerolog/log"
	"strconv"
)

type OrderService struct {
	db db.OrderDB
}

func NewOrderService(orderDB db.OrderDB) *OrderService {
	return &OrderService{db: orderDB}
}

func (s *OrderService) AddOrder(orderID, userID int) error {
	if !utils.ValidateLuhn(orderID) {
		log.Warn().Str("service", "Order service").Msg(fmt.Sprintf("order %d is not valid", orderID))
		return utils.ErrWrongOrderNumber
	}

	return s.db.AddOrder(strconv.Itoa(orderID), userID)
}

func (s *OrderService) GetOrders(userID int) ([]models.Order, error) {
	return s.db.GetOrders(userID)
}
