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

type UserHandler struct {
	log     *slog.Logger
	service UserService
	router  chi.Router
}

func NewUserHandler(log *slog.Logger, service UserService) *UserHandler {
	h := &UserHandler{
		log:     log,
		service: service,
		router:  chi.NewRouter(),
	}

	h.router.Post("/add_user", h.AddUserHandler)
	h.router.Get("/user_active_segments", h.UserActiveSegmentsHandler)

	return h
}

func (u UserHandler) Handler() http.Handler {
	return u.router
}

func (u *UserHandler) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request types.AddUserRequest
	if err := handlers.UnmarshalJSON(r, &request); err != nil {
		u.log.Error("failed to unmarshal request json: %v", err)
		handlers.RenderBadRequest(w, err)
		return
	}

	if err := u.service.AddUser(ctx, request); err != nil {
		u.log.Error(err.Error())
		handlers.RenderInternalError(w, err)
		return
	}

	handlers.RenderOK(w)
}

func (u *UserHandler) UserActiveSegmentsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request types.ActiveUserSegmentsRequest
	if err := handlers.UnmarshalJSON(r, &request); err != nil {
		u.log.Error("failed to unmarshal request json: %v", err)
		handlers.RenderBadRequest(w, err)
		return
	}

	response, err := u.service.UserActiveSegments(ctx, request)
	if err != nil {
		u.log.Error(err.Error())
		handlers.RenderBadRequest(w, err)
		return
	}

	handlers.RenderJSON(w, response)
}
