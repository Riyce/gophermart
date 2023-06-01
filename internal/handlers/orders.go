package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/riyce/gophermart/internal/utils"
	"github.com/rs/zerolog/log"
	"net/http"
)

const orderHandler string = "Order handler"

func (h *Handler) addOrder(ctx *gin.Context) {
	userID, ok := ctx.Get(ctxUser)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	intUserID, _ := userID.(int)
	log.Info().Str("service", orderHandler).Msg(fmt.Sprintf("user %d started addOrder", intUserID))
	var orderID int
	if err := ctx.BindJSON(&orderID); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := h.service.AddOrder(orderID, intUserID)
	if err != nil {
		if errors.Is(err, utils.ErrUsersOrderAlreadyExists) {
			ctx.Status(http.StatusOK)
		} else if errors.Is(err, utils.ErrOtherOrderAlreadyExists) {
			ctx.Status(http.StatusConflict)
		} else if errors.Is(err, utils.ErrWrongOrderNumber) {
			ctx.Status(http.StatusUnprocessableEntity)
		} else {
			ctx.Status(http.StatusInternalServerError)
		}
	} else {
		log.Info().Str("service", orderHandler).Msg(fmt.Sprintf("user %d addOrder OK", intUserID))
		ctx.Status(http.StatusAccepted)
	}
}

func (h *Handler) getOrders(ctx *gin.Context) {
	userID, ok := ctx.Get(ctxUser)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	intUserID, _ := userID.(int)
	log.Info().Str("service", orderHandler).Msg(fmt.Sprintf("user %d started getOrders", intUserID))
	orders, err := h.service.GetOrders(intUserID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Info().Str("service", orderHandler).Msg(fmt.Sprintf("user %d getOrders OK", intUserID))
	if len(orders) == 0 {
		ctx.JSON(http.StatusNoContent, nil)
	} else {
		ctx.JSON(http.StatusOK, orders)
	}
}
