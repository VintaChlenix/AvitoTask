package http

import (
	"context"
	"net/http"

	"avitoTask/internal/types"
	"avitoTask/pkg/handlers"
	"go.uber.org/zap"
)

type SegmentsService interface {
	CreateSegment(ctx context.Context, request types.CreateSegmentRequest) error
	DeleteSegment(ctx context.Context, request types.DeleteSegmentRequest) error
}

type SegmentsDelivery struct {
	log     *zap.SugaredLogger
	service SegmentsService
}

func NewSegmentsDelivery(log *zap.SugaredLogger, service SegmentsService) *SegmentsDelivery {
	return &SegmentsDelivery{
		log:     log,
		service: service,
	}
}

func (s *SegmentsDelivery) CreateSegmentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request types.CreateSegmentRequest
	if err := handlers.UnmarshalJSON(r, &request); err != nil {
		s.log.Errorf("failed to unmarshal request json: %v", err)
		handlers.RenderBadRequest(w, err)
		return
	}

	if err := s.service.CreateSegment(ctx, request); err != nil {
		s.log.Errorf(err.Error())
		handlers.RenderInternalError(w, err)
		return
	}

	handlers.RenderOK(w)
}

func (s *SegmentsDelivery) DeleteSegmentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request types.DeleteSegmentRequest
	if err := handlers.UnmarshalJSON(r, &request); err != nil {
		s.log.Errorf("failed to unmarshal request json: %v", err)
		handlers.RenderBadRequest(w, err)
		return
	}

	if err := s.service.DeleteSegment(ctx, request); err != nil {
		s.log.Errorf(err.Error())
		handlers.RenderInternalError(w, err)
		return
	}

	handlers.RenderOK(w)
}
