package app

import (
	"avitoTask/internal/db"
	"go.uber.org/zap"
)

type App struct {
	log *zap.SugaredLogger
	db  db.DB
}

func NewApp(log *zap.SugaredLogger, db db.DB) (*App, error) {
	return &App{
		log: log,
		db:  db,
	}, nil
}
