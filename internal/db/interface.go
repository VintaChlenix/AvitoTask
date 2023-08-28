package db

import (
	"avitoTask/internal/model"
)

type DB interface {
	CreateSegment(slug model.Slug) error
	DeleteSegment(slug model.Slug) error

	CreateUser(userID model.UserID, segmentsToAdd []model.Slug, segmentsToDelete []model.Slug) error
	SelectActiveSegments(userID model.UserID) ([]model.Slug, error)
}
