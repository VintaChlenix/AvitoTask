package http

import (
	"avitoTask/internal/types/dto"
	"avitoTask/pkg/handlers"
	"net/http"
)

func (a *app.App) CreateSegmentHandler(w http.ResponseWriter, r *http.Request) {
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

func (a *app.App) DeleteSegmentHandler(w http.ResponseWriter, r *http.Request) {
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
