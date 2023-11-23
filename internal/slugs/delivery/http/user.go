package http

import (
	"context"
	"log/slog"
	"net/http"

	"avitoTask/internal/slugs/types"
	"avitoTask/pkg/handlers"
	"github.com/go-chi/chi/v5"
)

type UserService interface {
	AddUser(ctx context.Context, request types.AddUserRequest) error
	UserActiveSegments(ctx context.Context, request types.ActiveUserSegmentsRequest) (*types.ActiveUserSegmentsResponse, error)
}

type User struct {
	service UserService
	router  chi.Router
}

func NewUser(service UserService) *User {
	h := &User{
		service: service,
		router:  chi.NewRouter(),
	}

	h.router.Post("/add_user", h.AddUserHandler)
	h.router.Get("/user_active_segments", h.UserActiveSegmentsHandler)

	return h
}

func (u User) Handler() http.Handler {
	return u.router
}

func (u *User) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request types.AddUserRequest
	if err := handlers.UnmarshalJSON(r, &request); err != nil {
		slog.Error("failed to unmarshal request json: %v", err)
		handlers.RenderBadRequest(w, err)
		return
	}

	if err := u.service.AddUser(ctx, request); err != nil {
		slog.Error(err.Error())
		handlers.RenderInternalError(w, err)
		return
	}

	handlers.RenderOK(w)
}

func (u *User) UserActiveSegmentsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request types.ActiveUserSegmentsRequest
	if err := handlers.UnmarshalJSON(r, &request); err != nil {
		slog.Error("failed to unmarshal request json: %v", err)
		handlers.RenderBadRequest(w, err)
		return
	}

	response, err := u.service.UserActiveSegments(ctx, request)
	if err != nil {
		slog.Error(err.Error())
		handlers.RenderBadRequest(w, err)
		return
	}

	handlers.RenderJSON(w, response)
}
