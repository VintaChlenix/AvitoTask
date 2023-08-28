package main

import (
	"avitoTask/internal/app"
	"github.com/go-chi/chi/v5"
)

func newRouter(a *app.App) chi.Router {
	r := chi.NewRouter()

	r.Post("/create", a.CreateSegmentHandler)
	r.Delete("/delete", a.DeleteSegmentHandler)
	r.Post("/add_user", a.AddUserHandler)
	r.Get("/user_active_segments", a.UserActiveSegmentsHandler)

	return r
}
