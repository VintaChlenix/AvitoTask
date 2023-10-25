package http

import (
	"context"
	"net/http"

	"avitoTask/internal/types"
	"avitoTask/pkg/handlers"
	"go.uber.org/zap"
)

type UsersService interface {
	AddUser(ctx context.Context, request types.AddUserRequest) error
	UserActiveSegments(ctx context.Context, request types.ActiveUserSegmentsRequest) (*types.ActiveUserSegmentsResponse, error)
}

type UsersDelivery struct {
	log     *zap.SugaredLogger
	service UsersService
}

func NewUsersDelivery(log *zap.SugaredLogger, service UsersService) *UsersDelivery {
	return &UsersDelivery{
		log:     log,
		service: service,
	}
}

func (u *UsersDelivery) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request types.AddUserRequest
	if err := handlers.UnmarshalJSON(r, &request); err != nil {
		u.log.Errorf("failed to unmarshal request json: %v", err)
		handlers.RenderBadRequest(w, err)
		return
	}

	if err := u.service.AddUser(ctx, request); err != nil {
		u.log.Errorf(err.Error())
		handlers.RenderInternalError(w, err)
		return
	}

	handlers.RenderOK(w)
}

func (u *UsersDelivery) UserActiveSegmentsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request types.ActiveUserSegmentsRequest
	if err := handlers.UnmarshalJSON(r, &request); err != nil {
		u.log.Errorf("failed to unmarshal request json: %v", err)
		handlers.RenderBadRequest(w, err)
		return
	}

	response, err := u.service.UserActiveSegments(ctx, request)
	if err != nil {
		u.log.Errorf(err.Error())
		handlers.RenderBadRequest(w, err)
		return
	}

	handlers.RenderJSON(w, response)
}
