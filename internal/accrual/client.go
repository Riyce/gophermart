package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/riyce/gophermart/internal/models"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strconv"
	"time"
)

type ClientDaemon struct {
	addr string
	db   DBOrder
}

const accrualDaemon string = "accrual daemon"

func NewClientDaemon(addr string, db DBOrder) *ClientDaemon {
	return &ClientDaemon{addr: addr, db: db}
}

func (d *ClientDaemon) sendRequest(orderID string) (order, int, error) {
	var order order

	requestURL := fmt.Sprintf("%s/api/orders/%s", d.addr, orderID)

	resp, err := http.Get(requestURL)

	if err != nil {
		return order, 0, err
	}

	body, _ := io.ReadAll(resp.Body)

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		err = json.Unmarshal(body, &order)
		if err != nil {
			return order, 0, err
		}
		return order, 0, nil
	case http.StatusTooManyRequests:
		seconds, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
		return order, seconds, errTooManyRequests
	case http.StatusNoContent:
		return order, 0, errNoContent
	case http.StatusInternalServerError:
		return order, 0, errServerError
	default:
		return order, 0, errWrongStatusCode
	}
}

func (d *ClientDaemon) RunDaemon() {
	for {
		orders, err := d.db.GetUnprocessedOrders()

		if err != nil {
			log.Error().
				Err(err).
				Str("service", accrualDaemon).
				Msg("error on get unprocessed orders from db")
		}
	loop:
		for _, order := range orders {
			log.Info().Str("service", accrualDaemon).
				Msg(fmt.Sprintf("Started process order %s %s", order.Order, order.Status))
			if order.Status == models.New {
				err := d.db.UpdateOrderStatus(order.Order)
				if err != nil {
					log.Error().
						Err(err).
						Str("service", accrualDaemon).
						Msg(fmt.Sprintf("error on update status PROCESSING for order %s", order.Order))
					continue loop
				}
			}

			order, secondsSleep, requestErr := d.sendRequest(order.Order)
			if requestErr != nil {
				if errors.Is(requestErr, errTooManyRequests) {
					log.Warn().
						Err(requestErr).
						Str("service", accrualDaemon).
						Msg("too many requests to accrual service")
					time.Sleep(time.Second * time.Duration(secondsSleep))
				} else {
					log.Warn().
						Err(requestErr).
						Msg("error on request to accrual service")
				}
				continue loop
			}

			err := d.db.ProcessOrder(order)
			if err != nil {
				log.Error().
					Err(err).
					Str("service", accrualDaemon).
					Msg(
						fmt.Sprintf("error on close PROCESSED/INVALID order %s with values: status - %s accrual - %f",
							order.Order, order.Status, order.Accrual),
					)
				continue loop
			}
		}
		time.Sleep(time.Second)
	}
}
