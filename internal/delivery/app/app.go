package app

import (
	"avitoTask/internal/repo"
	"avitoTask/internal/service"
	"go.uber.org/zap"
)

type App struct {
	log      *zap.SugaredLogger
	segments service.SegmentsService
	users    service.UsersService
}

func NewApp(log *zap.SugaredLogger, db repo.DB) (*App, error) {
	return &App{
		log: log,
		db:  db,
	}, nil
}
