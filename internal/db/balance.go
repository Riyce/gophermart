package db

import (
	"database/sql"
	"fmt"
	"github.com/riyce/gophermart/internal/models"
	"github.com/riyce/gophermart/internal/utils"
	"github.com/rs/zerolog/log"
)

const (
	withdrawsTableName string = "withdraws"
	balanceDB          string = "Balance/Withdraws DB"
)

type BalanceController struct {
	db *sql.DB
}

func NewBalanceController(db *sql.DB) *BalanceController {
	return &BalanceController{db: db}
}

func (b *BalanceController) GetBalance(userID int) (models.Balance, error) {
	var balance models.Balance
	query := fmt.Sprintf(getBalanceQuery, usersTableName)
	row := b.db.QueryRow(query, userID)
	if row.Err() != nil {
		log.Error().
			Err(row.Err()).
			Str("service", balanceDB).
			Msg(fmt.Sprintf("error on get balance for user %d", userID))
		return balance, utils.ErrSomethingWentWrong
	}

	if err := row.Scan(&balance.Withdrawn, &balance.Current); err != nil {
		log.Error().
			Err(err).
			Str("service", balanceDB).
			Msg(fmt.Sprintf("error on get balance for user %d", userID))
		return balance, utils.ErrSomethingWentWrong
	}

	return balance, nil
}

func (b *BalanceController) Withdraw(userID int, withdraw models.Withdraw) error {
	var balance models.Balance
	query := fmt.Sprintf(getBalanceQuery, usersTableName)
	row := b.db.QueryRow(query, userID)
	if row.Err() != nil {
		log.Error().
			Err(row.Err()).
			Str("service", balanceDB).
			Msg(fmt.Sprintf("error on withdraw balance for user %d", userID))
		return utils.ErrSomethingWentWrong
	}

	if err := row.Scan(&balance.Withdrawn, &balance.Current); err != nil {
		log.Error().
			Err(err).
			Str("service", balanceDB).
			Msg(fmt.Sprintf("error on withdraw balance for user %d", userID))
		return utils.ErrSomethingWentWrong
	}

	if balance.Current < withdraw.Sum {
		log.Warn().
			Str("service", balanceDB).
			Msg(fmt.Sprintf("user %d try to withdraw %f, current %f", userID, withdraw.Sum, balance.Current))
		return utils.ErrNotEnoughBalance
	}

	newBalance := balance.Current - withdraw.Sum
	newWithdraw := balance.Withdrawn + withdraw.Sum
	tx, txErr := b.db.Begin()
	if txErr != nil {
		log.Error().Err(txErr).Str("service", balanceDB).Msg("error on start transaction")
		return utils.ErrSomethingWentWrong
	}

	defer tx.Rollback()
	updateQuery := fmt.Sprintf(updateBalanceQuery, usersTableName)
	_, updateErr := tx.Exec(updateQuery, newBalance, newWithdraw, userID)
	if updateErr != nil {
		log.Error().
			Err(updateErr).
			Str("service", balanceDB).
			Msg(fmt.Sprintf("error on update user %d balance with current %f and withdrawn %f",
				userID, newBalance, newWithdraw))
		return utils.ErrSomethingWentWrong
	}

	createQuery := fmt.Sprintf(addWithdrawRecordQuery, withdrawsTableName)
	_, createErr := tx.Exec(createQuery, userID, withdraw.Order, withdraw.Sum)
	if createErr != nil {
		log.Error().
			Err(createErr).
			Str("service", balanceDB).
			Msg("error on add withdraw to history")
		return utils.ErrSomethingWentWrong
	}

	if commErr := tx.Commit(); commErr != nil {
		log.Error().
			Err(commErr).
			Str("service", balanceDB).
			Msg("error on commit changes")
		return utils.ErrSomethingWentWrong
	}

	return nil
}

func (b *BalanceController) GetWithdrawsHistory(userID int) ([]models.Withdraw, error) {
	var withdrawList []models.Withdraw
	query := fmt.Sprintf(getWithdrawsHistoryQuery, withdrawsTableName)
	rows, err := b.db.Query(query, userID)
	if err != nil || rows.Err() != nil {
		if err != nil {
			log.Error().
				Err(err).
				Str("service", balanceDB).
				Msg(fmt.Sprintf("error on get withdraws history for user %d", userID))
		} else {
			log.Error().
				Err(rows.Err()).
				Str("service", balanceDB).
				Msg(fmt.Sprintf("error on get withdraws history for user %d", userID))
		}
		return withdrawList, utils.ErrSomethingWentWrong
	}

	for rows.Next() {
		var record models.Withdraw
		if scanErr := rows.Scan(&record.Order, &record.Sum, &record.ProcessedAt); scanErr != nil {
			log.Error().
				Err(scanErr).
				Str("service", balanceDB).
				Msg(fmt.Sprintf("error on sczn withdraws history for user %d", userID))
			return withdrawList, utils.ErrSomethingWentWrong
		}
		withdrawList = append(withdrawList, record)
	}

	return withdrawList, nil
}
