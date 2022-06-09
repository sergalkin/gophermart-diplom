package config

import (
	"flag"
	"fmt"
	"sync"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

type config struct {
	ServerAddress  string `env:"RUN_ADDRESS" envDefault:"127.0.0.1:8081"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"http://127.0.0.1:8080"`
	DatabaseDSN    string `env:"DATABASE_URI" envDefault:"postgres://postgres:passwords4@localhost:5432/gophermart?sslmode=disable"`

	Logger *zap.Logger

	GinMode string `env:"GIN_MODE" envDefault:"debug"`
}

var (
	cfg  config
	once sync.Once
)

func NewConfig() *config {
	once.Do(func() {
		if err := env.Parse(&cfg); err != nil {
			fmt.Println(err.Error())
		}

		flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "listen address")
		flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "database dsn")
		flag.StringVar(&cfg.AccrualAddress, "r", cfg.AccrualAddress, "accrual service")
		flag.StringVar(&cfg.GinMode, "m", cfg.GinMode, "gin mode")

		flag.Parse()

		var z zap.Config
		if cfg.GinMode == "release" {
			z = zap.NewProductionConfig()
		} else {
			z = zap.NewDevelopmentConfig()
		}

		z.OutputPaths = []string{"stderr", "logs/zap.log"}
		l, err := z.Build()
		if err != nil {
			fmt.Println(err)
		}

		cfg.Logger = l
	})

	return &cfg
}

func Logger() *zap.Logger {
	return NewConfig().Logger
}

func ServerAddress() string {
	return NewConfig().ServerAddress
}

func AccrualAddress() string {
	return NewConfig().AccrualAddress
}

func DatabaseDSN() string {
	return NewConfig().DatabaseDSN
}

func GinMode() string {
	return NewConfig().GinMode
}
