package app

import (
	"avitoTask/internal/dto"
	"avitoTask/pkg/handlers"
	"fmt"
	"net/http"
)

func (a *App) CreateSegmentHandler(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateSegmentRequest
	if err := handlers.UnmarshalJSON(r, &request); err != nil {
		a.log.Errorf("failed to unmarshal request json: %v", err)
		handlers.RenderBadRequest(w, err)
		return
	}

	if err := a.createSegmentHandler(request); err != nil {
		a.log.Errorf(err.Error())
		handlers.RenderInternalError(w, err)
		return
	}

	handlers.RenderOK(w)
}

func (a *App) createSegmentHandler(request dto.CreateSegmentRequest) error {
	if err := a.db.CreateSegment(request.Slug); err != nil {
		return fmt.Errorf("failed to create segment: %w", err)
	}
	return nil
}

func (a *App) DeleteSegmentHandler(w http.ResponseWriter, r *http.Request) {
	var request dto.DeleteSegmentRequest
	if err := handlers.UnmarshalJSON(r, &request); err != nil {
		a.log.Errorf("failed to unmarshal request json: %v", err)
		handlers.RenderBadRequest(w, err)
		return
	}

	if err := a.deleteSegmentHandler(request); err != nil {
		a.log.Errorf(err.Error())
		handlers.RenderInternalError(w, err)
		return
	}

	handlers.RenderOK(w)
}

func (a *App) deleteSegmentHandler(request dto.DeleteSegmentRequest) error {
	if err := a.db.DeleteSegment(request.Slug); err != nil {
		return fmt.Errorf("failed to delete segment: %w", err)
	}
	return nil
}
