package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"avitoTask/internal/config"
	handlers "avitoTask/internal/slugs/delivery/http"
	"avitoTask/internal/slugs/repo/postgres"
	"avitoTask/internal/slugs/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.GetConfig(config.PATH)
	if err != nil {
		return err
	}

	dbClient, err := pgxpool.New(context.TODO(), cfg.Database.ConnectionString)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer dbClient.Close()

	segmentDelivery := handlers.NewSegment(service.NewSegment(postgres.NewSegment(dbClient)))
	userDelivery := handlers.NewUser(service.NewUser(postgres.NewUser(dbClient)))
	slog.Info("app initialized")

	http.Handle("/api/segment/", segmentDelivery.Handler())
	http.Handle("/api/user/", userDelivery.Handler())

	srv := http.Server{
		Addr: cfg.Server.URL,
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		slog.Info("Starting server on: %s", cfg.Server.URL)
		lerr := srv.ListenAndServe()
		if errors.Is(lerr, http.ErrServerClosed) {
			return nil
		}

		return lerr
	})

	eg.Go(func() error {
		<-ctx.Done()
		return srv.Shutdown(ctx)
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	slog.Info("Shutdown app")
	return nil
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	if err := run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
