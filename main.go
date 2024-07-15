package main

import (
	"auth-api/config"
	"auth-api/internal/manager"
	"auth-api/logger"
	"go.uber.org/fx"
	"log"
)

func main() {
	app := fx.New(
		fx.Provide(
			config.InitConfig,
			logger.InitSlogLogger,
			manager.InitManager,
		),
		fx.Invoke(manager.Run),
	)

	app.Run()
	if err := app.Err(); err != nil {
		log.Fatal(err)
	}
}
