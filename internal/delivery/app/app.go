package app

import (
	delivery "avitoTask/internal/delivery/handlers/http"
)

type App struct {
	segments *delivery.SegmentsDelivery
	users    *delivery.UsersDelivery
}

func NewApp(segments *delivery.SegmentsDelivery, users *delivery.UsersDelivery) (*App, error) {
	return &App{
		segments: segments,
		users:    users,
	}, nil
}
