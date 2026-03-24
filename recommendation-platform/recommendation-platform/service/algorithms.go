package service

import (
	"sort"

	"recommendation-platform/model"
	"recommendation-platform/pipeline"
	"recommendation-platform/repository"
)

// Algorithm is a small wrapper so changing algorithms only touches this file.
type Algorithm interface {
	Name() string
	Recommend(userID string, limit int) []model.Video
}

type AlgorithmRegistry struct {
	items map[string]Algorithm
}

func NewAlgorithmRegistry() *AlgorithmRegistry {
	r := &AlgorithmRegistry{items: map[string]Algorithm{}}
	r.Register(&DefaultAlgorithm{})
	r.Register(&LatestAlgorithm{})
	r.Register(&HotAlgorithm{})
	r.Register(&FollowingAlgorithm{})
	return r
}

func (r *AlgorithmRegistry) Register(a Algorithm) {
	if a == nil || a.Name() == "" {
		return
	}
	r.items[a.Name()] = a
}

func (r *AlgorithmRegistry) Get(name string) Algorithm {
	if a, ok := r.items[name]; ok {
		return a
	}
	return r.items["default"]
}

type DefaultAlgorithm struct{}

func (a *DefaultAlgorithm) Name() string { return "default" }

type DefaultAlgoConfig struct {
	MaxAgeDays               int
	TopTagsLimit             int
	GraphBoostPerEdge        int
	RecencyMaxDays           int
	AuthorDiversityPenalty   int
}

func defaultAlgoConfig() DefaultAlgoConfig {
	return DefaultAlgoConfig{
		MaxAgeDays:             30,
		TopTagsLimit:           3,
		GraphBoostPerEdge:      3,
		RecencyMaxDays:         30,
		AuthorDiversityPenalty: 2,
	}
}

func (a *DefaultAlgorithm) Recommend(userID string, limit int) []model.Video {
	cfg := defaultAlgoConfig()
	engine := pipeline.Engine{
		Sources: []pipeline.Source{
			&pipeline.InNetworkSource{},
			&pipeline.TagSearchSource{},
			&pipeline.GraphSource{},
			&pipeline.IndexSource{},
		},
		Filters: []pipeline.Filter{
			&pipeline.PreviouslySeenFilter{},
			&pipeline.DedupFilter{},
			&pipeline.AgeFilter{MaxAgeDays: cfg.MaxAgeDays},
		},
		Scorers: []pipeline.Scorer{
			&pipeline.ViewsScorer{},
			&pipeline.BehaviorScorer{},
			&pipeline.TagSimilarityScorer{},
			&pipeline.GraphProximityScorer{BoostPerEdge: cfg.GraphBoostPerEdge},
			&pipeline.RecencyScorer{MaxDays: cfg.RecencyMaxDays},
			&pipeline.AuthorDiversityScorer{PenaltyPerExtra: cfg.AuthorDiversityPenalty},
		},
		Selector: &pipeline.TopKSelector{},
		K:        limit,
	}

	indexVideos := pipeline.LoadIndexVideos(pipeline.Query{})
	indexMap := map[string]model.Video{}
	for _, v := range indexVideos {
		indexMap[v.ID] = v
	}

	behaviors := repository.GetBehaviors()
	follows := repository.GetFollows()
	topTags := topTagsFromBehaviors(behaviors, indexMap, cfg.TopTagsLimit)

	graph := pipeline.BuildVideoGraph(indexMap)

	result := engine.Run(pipeline.Query{
		UserID:      userID,
		IndexVideos: indexMap,
		Behaviors:   behaviors,
		TopTags:     topTags,
		Follows:     follows,
		Graph:       graph,
	})
	out := make([]model.Video, 0, len(result.Selected))
	for _, c := range result.Selected {
		out = append(out, c.Video)
	}
	return out
}

type LatestAlgorithm struct{}

func (a *LatestAlgorithm) Name() string { return "latest" }

func (a *LatestAlgorithm) Recommend(userID string, limit int) []model.Video {
	videos := pipeline.LoadIndexVideos(pipeline.Query{})
	sort.Slice(videos, func(i, j int) bool {
		return videos[i].CreatedAt > videos[j].CreatedAt
	})
	if len(videos) > limit {
		return videos[:limit]
	}
	return videos
}

type HotAlgorithm struct{}

func (a *HotAlgorithm) Name() string { return "hot" }

func (a *HotAlgorithm) Recommend(userID string, limit int) []model.Video {
	videos := pipeline.LoadIndexVideos(pipeline.Query{})
	behaviors := repository.GetBehaviors()
	counts := map[string]int{}
	for _, b := range behaviors {
		switch b.Type {
		case "like":
			counts[b.VideoID] += 3
		case "share":
			counts[b.VideoID] += 5
		case "watch":
			counts[b.VideoID] += 1
		}
	}
	sort.Slice(videos, func(i, j int) bool {
		return counts[videos[i].ID] > counts[videos[j].ID]
	})
	if len(videos) > limit {
		return videos[:limit]
	}
	return videos
}

type FollowingAlgorithm struct{}

func (a *FollowingAlgorithm) Name() string { return "following" }

func (a *FollowingAlgorithm) Recommend(userID string, limit int) []model.Video {
	if userID == "" {
		return []model.Video{}
	}
	follows := repository.GetFollows()
	following := map[string]bool{}
	for _, f := range follows {
		if f.UserID == userID && f.Active {
			following[f.AuthorID] = true
		}
	}
	if len(following) == 0 {
		return []model.Video{}
	}
	videos := pipeline.LoadIndexVideos(pipeline.Query{})
	out := make([]model.Video, 0)
	for _, v := range videos {
		if following[v.AuthorID] {
			out = append(out, v)
		}
	}
	if len(out) > limit {
		return out[:limit]
	}
	return out
}
