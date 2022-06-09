package server

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/sergalkin/gophermart-diplom/internal/gophermart/config"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/routes"
	"github.com/sergalkin/gophermart-diplom/internal/gophermart/storage"
)

func NewServer() error {
	if err := setupLogging(); err != nil {
		return err
	}

	gin.SetMode(config.GinMode())
	e := gin.Default()

	loadGlobalMiddlewares(e)
	loadRoutes(e, storage.NewStorage())

	if err := e.Run(config.ServerAddress()); err != nil {
		return fmt.Errorf("error in starting gin server %w", err)
	}

	return nil
}

func setupLogging() error {
	gin.DisableConsoleColor()

	fg, err := os.Create("logs/gin.log")
	if err != nil {
		return fmt.Errorf("error in creating gin.log file; %w", err)
	}
	gin.DefaultWriter = fg

	fe, fErr := os.Create("logs/error.log")
	if fErr != nil {
		return fmt.Errorf("error in creating error.log file; %w", fErr)
	}
	gin.DefaultErrorWriter = fe

	return nil
}

func loadGlobalMiddlewares(e *gin.Engine) {
	//todo gzip
}

func loadRoutes(e *gin.Engine, s storage.Storage) {
	routes.NewAPIRoutes(e, s).Setup()
}
