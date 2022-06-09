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

	api := a.engine.Group("/api", middleware.JSONConcern())
	{
		users := api.Group("/user")
		{
			auth := handlers.NewAuthHandler(service.NewAuthService(a.storage))

			users.POST("/register", auth.Register())
			users.POST("/login", auth.Login())
		}
	}
}
