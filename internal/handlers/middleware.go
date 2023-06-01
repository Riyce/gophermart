package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/riyce/gophermart/internal/utils"
	"net/http"
)

const ctxUser string = "userID"

func (h *Handler) userIdentity(ctx *gin.Context) {
	key := ctx.GetHeader(utils.AuthorizationHeader)
	if key == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userID, err := h.service.Auth.GetUserID(key)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Set(ctxUser, userID)
}
