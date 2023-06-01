package handlers

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/riyce/gophermart/internal/services"
)

type Handler struct {
	service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes(debug bool) *gin.Engine {
	router := gin.New()

	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router.Use(gzip.Gzip(gzip.BestCompression))

	api := router.Group("/api/user")
	{
		api.POST("/register", h.signUp)
		api.POST("/login", h.signIn)
		orders := api.Group("/orders", h.userIdentity)
		{
			orders.POST("/", h.addOrder)
			orders.GET("/", h.getOrders)
		}
		balance := api.Group("/balance", h.userIdentity)
		{
			balance.GET("/", h.getBalance)
			balance.POST("/withdraw", h.withdraw)
		}
		withdrawals := api.Group("/withdrawals", h.userIdentity)
		{
			withdrawals.GET("/", h.getWithdrawsHistory)
		}
	}
	return router
}
