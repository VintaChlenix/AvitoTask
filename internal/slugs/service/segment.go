package service

import (
	"context"
	"fmt"

	"avitoTask/internal/slugs/types"
)

type SegmentRepo interface {
	CreateSegment(ctx context.Context, slug types.Slug) error
	DeleteSegment(ctx context.Context, slug types.Slug) error
}

type SegmentStore struct {
	repo SegmentRepo
}

func NewSegmentsService(repo SegmentRepo) *SegmentStore {
	return &SegmentStore{repo: repo}
}

func (s *SegmentStore) CreateSegment(ctx context.Context, request types.CreateSegmentRequest) error {
	if err := s.repo.CreateSegment(ctx, request.Slug); err != nil {
		return fmt.Errorf("failed to create segment: %w", err)
	}
	return nil
}

func (s *SegmentStore) DeleteSegment(ctx context.Context, request types.DeleteSegmentRequest) error {
	if err := s.repo.DeleteSegment(ctx, request.Slug); err != nil {
		return fmt.Errorf("failed to delete segment: %w", err)
	}
	return nil
}
