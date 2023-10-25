package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"avitoTask/internal/delivery/app"
	delivery "avitoTask/internal/delivery/handlers/http"
	"avitoTask/internal/repo/postgres"
	"avitoTask/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func run(log *zap.SugaredLogger) error {
	config, err := GetConfig()
	if err != nil {
		log.Panic(err.Error())
	}

	dbClient, err := pgxpool.New(context.TODO(), config.Database.ConnectionString)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer dbClient.Close()

	segmentsDelivery := delivery.NewSegmentsDelivery(log, service.NewSegmentsService(postgres.NewSegmentsRepo(dbClient)))
	usersDelivery := delivery.NewUsersDelivery(log, service.NewUsersService(postgres.NewUsersRepo(dbClient)))

	avitoApp, err := app.NewApp(segmentsDelivery, usersDelivery)
	if err != nil {
		return fmt.Errorf("failed to initialize app: %v", err)
	}
	log.Infoln("app initialized")
	log.Infof("Starting server on: %s", config.Server.URL)

	srv := http.Server{}
	srv.Addr = config.Server.URL
	srv.Handler = newRouter(avitoApp)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve returned err: %v", err)
		}
	}()

	<-ctx.Done()
	log.Infoln("got interruption signal")
	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Infof("server shutdown returned an err: %v\n", err)
	}

	log.Infoln("final")
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
