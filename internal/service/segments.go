package service

import (
	"context"
	"fmt"

	"avitoTask/internal/types"
)

type SegmentsRepo interface {
	CreateSegment(ctx context.Context, slug types.Slug) error
	DeleteSegment(ctx context.Context, slug types.Slug) error
}

type SegmentsService struct {
	repo SegmentsRepo
}

func NewSegmentsService(repo SegmentsRepo) *SegmentsService {
	return &SegmentsService{repo: repo}
}

func (s *SegmentsService) CreateSegment(ctx context.Context, request types.CreateSegmentRequest) error {
	if err := s.repo.CreateSegment(ctx, request.Slug); err != nil {
		return fmt.Errorf("failed to create segment: %w", err)
	}
	return nil
}

func (s *SegmentsService) DeleteSegment(ctx context.Context, request types.DeleteSegmentRequest) error {
	if err := s.repo.DeleteSegment(ctx, request.Slug); err != nil {
		return fmt.Errorf("failed to delete segment: %w", err)
	}
	return nil
}
