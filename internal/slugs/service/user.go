package service

import (
	"context"
	"fmt"

	types2 "avitoTask/internal/slugs/types"
)

type UserRepo interface {
	CreateUser(ctx context.Context, userID types2.UserID, segmentsToAdd []types2.Slug, segmentsToDelete []types2.Slug) error
	SelectActiveSegments(ctx context.Context, userID types2.UserID) ([]types2.Slug, error)
}

type UserStore struct {
	repo UserRepo
}

func NewUsersService(repo UserRepo) *UserStore {
	return &UserStore{repo: repo}
}

func (u *UserStore) AddUser(ctx context.Context, request types2.AddUserRequest) error {
	if err := u.repo.CreateUser(ctx, request.UserID, request.SegmentsToAdd, request.SegmentsToDelete); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (u *UserStore) UserActiveSegments(ctx context.Context, request types2.ActiveUserSegmentsRequest) (*types2.ActiveUserSegmentsResponse, error) {
	activeSegments, err := u.repo.SelectActiveSegments(ctx, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to select active segments: %w", err)
	}
	return &types2.ActiveUserSegmentsResponse{ActiveSegments: activeSegments}, nil
}
