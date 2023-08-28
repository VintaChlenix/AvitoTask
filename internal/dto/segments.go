package dto

import "avitoTask/internal/model"

type CreateSegmentRequest struct {
	Slug model.Slug `json:"slug"`
}

type DeleteSegmentRequest struct {
	Slug model.Slug `json:"slug"`
}
