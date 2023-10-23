package main

import (
	"avitoTask/internal/delivery/app"
	"avitoTask/internal/repo/postgres"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
)

func run(log *zap.SugaredLogger) error {
	config, err := GetConfig()
	if err != nil {
		log.Panic(err.Error())
	}

	dbClient, err := postgres.NewClient(config.Database.ConnectionString)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	avitoApp, err := app.NewApp(log, dbClient)
	if err != nil {
		return fmt.Errorf("failed to initialize app: %v", err)
	}
	log.Infoln("app initialized")
	log.Infof("Starting server on: %s", config.Server.URL)

	if err := http.ListenAndServe(config.Server.URL, newRouter(avitoApp)); err != nil {
		return fmt.Errorf("failed to run server: %w", err)
	}

	return nil
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	if err := run(sugar); err != nil {
		sugar.Error(err.Error())
		os.Exit(1)
	}
}
