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

const authHandler string = "Auth handler"

func (h *Handler) signIn(ctx *gin.Context) {
	var user models.User
	if err := ctx.BindJSON(&user); err != nil {
		log.Error().Err(err).Str("service", authHandler).Msg("error on bind json")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	getErr := h.service.Auth.GetUser(&user)
	if getErr != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	log.Info().Str("service", authHandler).Msg(fmt.Sprintf("user %s logged in", user.Login))
	ctx.Header(utils.AuthorizationHeader, user.APIKey)
	ctx.JSON(http.StatusOK, nil)
}

func (h *Handler) signUp(ctx *gin.Context) {
	var user models.User
	if err := ctx.BindJSON(&user); err != nil {
		log.Error().Err(err).Str("service", authHandler).Msg("error on bind json")
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	log.Info().Str("service", authHandler).Msg(fmt.Sprintf("try to create user %s", user.Login))
	createErr := h.service.Auth.CreateUser(&user)
	if createErr != nil {
		if errors.Is(createErr, utils.ErrUserAlreadyExists) {
			ctx.AbortWithStatus(http.StatusConflict)
			return
		} else {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	log.Info().Str("service", authHandler).Msg(fmt.Sprintf("user %s created", user.Login))
	ctx.Header(utils.AuthorizationHeader, user.APIKey)
	ctx.JSON(http.StatusOK, nil)
}
