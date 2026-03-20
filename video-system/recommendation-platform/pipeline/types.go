package pipeline

import "recommendation-platform/model"

type Query struct {
	UserID      string
	IndexVideos map[string]model.Video
	Behaviors   []model.Behavior
	TopTags     []string
	Follows     []model.Follow
	Graph       *VideoGraph
}

type Candidate struct {
	Video model.Video
	Score int
}
