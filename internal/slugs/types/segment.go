package types

type Slug string

type CreateSegmentRequest struct {
	Slug Slug `json:"slug"`
}

type DeleteSegmentRequest struct {
	Slug Slug `json:"slug"`
}
