package app

import (
	"avitoTask/internal/dto"
	"avitoTask/pkg/handlers"
	"context"
	"fmt"
	"net/http"
)

func (a *App) CreateSegmentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request dto.CreateSegmentRequest
	if err := handlers.UnmarshalJSON(r, &request); err != nil {
		a.log.Errorf("failed to unmarshal request json: %v", err)
		handlers.RenderBadRequest(w, err)
		return
	}

	if err := a.createSegmentHandler(ctx, request); err != nil {
		a.log.Errorf(err.Error())
		handlers.RenderInternalError(w, err)
		return
	}

	handlers.RenderOK(w)
}

func (a *App) createSegmentHandler(ctx context.Context, request dto.CreateSegmentRequest) error {
	if err := a.db.CreateSegment(ctx, request.Slug); err != nil {
		return fmt.Errorf("failed to create segment: %w", err)
	}
	return nil
}

func (a *App) DeleteSegmentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request dto.DeleteSegmentRequest
	if err := handlers.UnmarshalJSON(r, &request); err != nil {
		a.log.Errorf("failed to unmarshal request json: %v", err)
		handlers.RenderBadRequest(w, err)
		return
	}

	if err := a.deleteSegmentHandler(ctx, request); err != nil {
		a.log.Errorf(err.Error())
		handlers.RenderInternalError(w, err)
		return
	}

	handlers.RenderOK(w)
}

func (a *App) deleteSegmentHandler(ctx context.Context, request dto.DeleteSegmentRequest) error {
	if err := a.db.DeleteSegment(ctx, request.Slug); err != nil {
		return fmt.Errorf("failed to delete segment: %w", err)
	}
	return nil
}
