package db

import (
	"avitoTask/internal/model"
	"context"
)

type DB interface {
	CreateSegment(ctx context.Context, slug model.Slug) error
	DeleteSegment(ctx context.Context, slug model.Slug) error

	CreateUser(ctx context.Context, userID model.UserID, segmentsToAdd []model.Slug, segmentsToDelete []model.Slug) error
	SelectActiveSegments(ctx context.Context, userID model.UserID) ([]model.Slug, error)
}
