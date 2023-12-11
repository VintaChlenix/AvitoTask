package types

type UserID int

type AddUserRequest struct {
	UserID           UserID `json:"user_id"`
	SegmentsToAdd    []Slug `json:"segments_to_add"`
	SegmentsToDelete []Slug `json:"segments_to_delete"`
}

type ActiveUserSegmentsRequest struct {
	UserID UserID `json:"user_id"`
}

type ActiveUserSegmentsResponse struct {
	ActiveSegments []Slug `json:"active_segments"`
}
