package services

import (
	"fmt"
	"github.com/riyce/gophermart/internal/db"
	"github.com/riyce/gophermart/internal/models"
	"github.com/riyce/gophermart/internal/utils"
	"github.com/rs/zerolog/log"
	"strconv"
)

type BalanceService struct {
	db db.BalanceDB
}

const balanceService string = "Balance service"

func NewBalanceService(balanceDB db.BalanceDB) *BalanceService {
	return &BalanceService{db: balanceDB}
}

func (o *BalanceService) GetBalance(userID int) (models.Balance, error) {
	return o.db.GetBalance(userID)
}

func (o *BalanceService) Withdraw(userID int, withdraw models.Withdraw) error {
	intOrderID, err := strconv.Atoi(withdraw.Order)
	if err != nil {
		log.Error().
			Err(err).
			Str("service", balanceService).
			Msg(fmt.Sprintf("can't convert order %s to int", withdraw.Order))
		return utils.ErrWrongOrderNumber
	}

	if !utils.ValidateLuhn(intOrderID) {
		log.Warn().
			Str("service", balanceService).
			Msg(fmt.Sprintf("order %s is not valid", withdraw.Order))
		return utils.ErrWrongOrderNumber
	}

	return o.db.Withdraw(userID, withdraw)
}

func (o *BalanceService) GetWithdrawsHistory(userID int) ([]models.Withdraw, error) {
	return o.db.GetWithdrawsHistory(userID)
}
