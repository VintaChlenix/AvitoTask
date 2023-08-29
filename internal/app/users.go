package app

import (
	"avitoTask/internal/dto"
	"avitoTask/pkg/handlers"
	"context"
	"fmt"
	"net/http"
)

func (a *App) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request dto.AddUserRequest
	if err := handlers.UnmarshalJSON(r, &request); err != nil {
		a.log.Errorf("failed to unmarshal request json: %v", err)
		handlers.RenderBadRequest(w, err)
		return
	}

	if err := a.addUserHandler(ctx, request); err != nil {
		a.log.Errorf(err.Error())
		handlers.RenderInternalError(w, err)
		return
	}

	handlers.RenderOK(w)
}

func (a *App) addUserHandler(ctx context.Context, request dto.AddUserRequest) error {
	if err := a.db.CreateUser(ctx, request.UserID, request.SegmentsToAdd, request.SegmentsToDelete); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (a *App) UserActiveSegmentsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request dto.ActiveSegmentsRequest
	if err := handlers.UnmarshalJSON(r, &request); err != nil {
		a.log.Errorf("failed to unmarshal request json: %v", err)
		handlers.RenderBadRequest(w, err)
		return
	}

	response, err := a.userActiveSegmentsHandler(ctx, request)
	if err != nil {
		a.log.Errorf(err.Error())
		handlers.RenderBadRequest(w, err)
		return
	}

	handlers.RenderJSON(w, response)
}

func (a *App) userActiveSegmentsHandler(ctx context.Context, request dto.ActiveSegmentsRequest) (*dto.ActiveSegmentsResponse, error) {
	activeSegments, err := a.db.SelectActiveSegments(ctx, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to select active segments: %w", err)
	}
	return &dto.ActiveSegmentsResponse{ActiveSegments: activeSegments}, nil
}
