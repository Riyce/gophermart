package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/riyce/gophermart/internal/models"
	"github.com/riyce/gophermart/internal/utils"
	"github.com/rs/zerolog/log"
	"net/http"
)

const balanceHandler string = "Balance handler"

func (h *Handler) getBalance(ctx *gin.Context) {
	userID, ok := ctx.Get(ctxUser)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	intUserID, _ := userID.(int)
	log.Info().Str("service", balanceHandler).Msg(fmt.Sprintf("user %d started getBalance", intUserID))
	balance, err := h.service.Balance.GetBalance(intUserID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Info().Str("service", balanceHandler).Msg(fmt.Sprintf("user %d getBalance OK", intUserID))
	ctx.JSON(http.StatusOK, balance)
}

func (h *Handler) withdraw(ctx *gin.Context) {
	userID, ok := ctx.Get(ctxUser)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	intUserID, _ := userID.(int)
	log.Info().Str("service", balanceHandler).
		Msg(fmt.Sprintf("user %d started withdraw", intUserID))
	var withdraw models.Withdraw
	if err := ctx.BindJSON(&withdraw); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := h.service.Balance.Withdraw(intUserID, withdraw)
	if err != nil {
		if errors.Is(err, utils.ErrNotEnoughBalance) {
			ctx.AbortWithStatus(http.StatusPaymentRequired)
			return
		} else if errors.Is(err, utils.ErrWrongOrderNumber) {
			ctx.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		} else {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	log.Info().Str("service", balanceHandler).
		Msg(fmt.Sprintf("user %d withdraw OK: withdrawed %f for order %s",
			intUserID, withdraw.Sum, withdraw.Order))
	ctx.Status(http.StatusOK)
}

func (h *Handler) getWithdrawsHistory(ctx *gin.Context) {
	userID, ok := ctx.Get(ctxUser)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	intUserID, _ := userID.(int)
	log.Info().Str("service", balanceHandler).
		Msg(fmt.Sprintf("user %d started getWithdrawsHistory", intUserID))
	withdraws, err := h.service.Balance.GetWithdrawsHistory(intUserID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Info().Str("service", balanceHandler).
		Msg(fmt.Sprintf("user %d: getWithdrawsHistory OK", intUserID))
	if len(withdraws) == 0 {
		ctx.JSON(http.StatusNoContent, nil)
	} else {
		ctx.JSON(http.StatusOK, withdraws)
	}
}
