package service

import (
	"context"
	"fmt"

	"avitoTask/internal/slugs/types"
)

type UserRepo interface {
	CreateUser(ctx context.Context, userID types.UserID, segmentsToAdd []types.Slug, segmentsToDelete []types.Slug) error
	SelectActiveSegments(ctx context.Context, userID types.UserID) ([]types.Slug, error)
}

type User struct {
	repo UserRepo
}

func NewUser(repo UserRepo) *User {
	return &User{repo: repo}
}

func (u *User) AddUser(ctx context.Context, request types.AddUserRequest) error {
	if err := u.repo.CreateUser(ctx, request.UserID, request.SegmentsToAdd, request.SegmentsToDelete); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (u *User) UserActiveSegments(ctx context.Context, request types.ActiveUserSegmentsRequest) (*types.ActiveUserSegmentsResponse, error) {
	activeSegments, err := u.repo.SelectActiveSegments(ctx, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to select active segments: %w", err)
	}
	return &types.ActiveUserSegmentsResponse{ActiveSegments: activeSegments}, nil
}
