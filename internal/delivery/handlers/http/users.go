package http

import (
	"avitoTask/internal/types"
	"avitoTask/pkg/handlers"
	"net/http"
)

func (a *app.App) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request types.AddUserRequest
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

func (a *app.App) UserActiveSegmentsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request types.ActiveSegmentsRequest
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
