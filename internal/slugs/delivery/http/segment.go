package http

import (
	"context"
	"log/slog"
	"net/http"

	"avitoTask/internal/slugs/types"
	"avitoTask/pkg/handlers"
	"github.com/go-chi/chi/v5"
)

type SegmentService interface {
	CreateSegment(ctx context.Context, request types.CreateSegmentRequest) error
	DeleteSegment(ctx context.Context, request types.DeleteSegmentRequest) error
}

type Segment struct {
	service SegmentService
	router  chi.Router
}

func NewSegment(service SegmentService) *Segment {
	h := &Segment{
		service: service,
		router:  chi.NewRouter(),
	}

	h.router.Post("/create", h.CreateSegmentHandler)
	h.router.Delete("/delete", h.DeleteSegmentHandler)

	return h
}

func (s *Segment) Handler() http.Handler {
	return s.router
}

func (s *Segment) CreateSegmentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request types.CreateSegmentRequest
	if err := handlers.UnmarshalJSON(r, &request); err != nil {
		slog.Error("failed to unmarshal request json: %v", err)
		handlers.RenderBadRequest(w, err)
		return
	}

	if err := s.service.CreateSegment(ctx, request); err != nil {
		slog.Error(err.Error())
		handlers.RenderInternalError(w, err)
		return
	}

	handlers.RenderOK(w)
}

func (s *Segment) DeleteSegmentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request types.DeleteSegmentRequest
	if err := handlers.UnmarshalJSON(r, &request); err != nil {
		slog.Error("failed to unmarshal request json: %v", err)
		handlers.RenderBadRequest(w, err)
		return
	}

	if err := s.service.DeleteSegment(ctx, request); err != nil {
		slog.Error(err.Error())
		handlers.RenderInternalError(w, err)
		return
	}

	handlers.RenderOK(w)
}
