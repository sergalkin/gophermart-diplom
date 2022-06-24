package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/sergalkin/gophermart-diplom/internal/gophermart/handlers"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/middleware"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/service"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/storage"
)

type apiRoutes struct {
	engine  *gin.Engine
	storage storage.Storage
}

func NewAPIRoutes(e *gin.Engine, s storage.Storage) *apiRoutes {
	return &apiRoutes{
		engine:  e,
		storage: s,
	}
}

func (a *apiRoutes) Setup() {

	api := a.engine.Group("/api")
	{
		user := api.Group("/user")
		{
			auth := handlers.NewAuthHandler(service.NewAuthService(a.storage))

			user.POST("/register", auth.Register(), middleware.JSONConcern())
			user.POST("/login", auth.Login(), middleware.JSONConcern())

			orders := user.Group("/orders", middleware.AuthConcern())
			{
				ordersHandler := handlers.NewOrdersHandler(service.NewOrdersService(a.storage))
				orders.POST("/", ordersHandler.Post(), middleware.TextPlainConcern())
				orders.GET("/", ordersHandler.Get())
			}

			balance := user.Group("/balance", middleware.AuthConcern())
			{
				balanceHandler := handlers.NewBalanceHandler(service.NewBalanceService(a.storage))
				balance.GET("/", balanceHandler.Check())
				balance.POST("/withdraw", balanceHandler.Withdraw(), middleware.JSONConcern())
				balance.GET("/withdrawals", balanceHandler.Withdrawals())
			}
		}
	}
}
