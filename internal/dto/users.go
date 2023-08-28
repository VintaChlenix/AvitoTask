package dto

import "avitoTask/internal/model"

type AddUserRequest struct {
	UserID           model.UserID `json:"user_id"`
	SegmentsToAdd    []model.Slug `json:"segments_to_add"`
	SegmentsToDelete []model.Slug `json:"segments_to_delete"`
}

type ActiveSegmentsRequest struct {
	UserID model.UserID `json:"user_id"`
}

type ActiveSegmentsResponse struct {
	ActiveSegments []model.Slug `json:"active_segments"`
}
