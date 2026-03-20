package service

import (
	"sort"

	"recommendation-platform/model"
	"recommendation-platform/pipeline"
	"recommendation-platform/repository"
)

type RecommendService struct {
}

func NewRecommendService() *RecommendService {

	return &RecommendService{}

}

func (s *RecommendService) Recommend(userID string, recType string, limit int) []model.Video {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	switch recType {
	case "latest":
		return latestVideos(limit)
	case "hot":
		return hotVideos(limit)
	case "following":
		return followingVideos(userID, limit)
	default:
		return s.recommendDefault(userID, limit)
	}
}

func (s *RecommendService) recommendDefault(userID string, limit int) []model.Video {

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
			&pipeline.AgeFilter{MaxAgeDays: 30},
		},
		Scorers: []pipeline.Scorer{
			&pipeline.ViewsScorer{},
			&pipeline.BehaviorScorer{},
			&pipeline.TagSimilarityScorer{},
			&pipeline.GraphProximityScorer{BoostPerEdge: 3},
			&pipeline.RecencyScorer{MaxDays: 30},
			&pipeline.AuthorDiversityScorer{PenaltyPerExtra: 2},
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
	topTags := topTagsFromBehaviors(behaviors, indexMap, 3)

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

func latestVideos(limit int) []model.Video {
	videos := pipeline.LoadIndexVideos(pipeline.Query{})
	sort.Slice(videos, func(i, j int) bool {
		return videos[i].CreatedAt > videos[j].CreatedAt
	})
	if len(videos) > limit {
		return videos[:limit]
	}
	return videos
}

func hotVideos(limit int) []model.Video {
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

func followingVideos(userID string, limit int) []model.Video {
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

func behaviorWeight(t string) int {
	switch t {
	case "like":
		return 5
	case "share":
		return 8
	case "watch":
		return 1
	default:
		return 0
	}
}

func topTagsFromBehaviors(behaviors []model.Behavior, videos map[string]model.Video, limit int) []string {
	counts := map[string]int{}
	for _, b := range behaviors {
		v, ok := videos[b.VideoID]
		if !ok {
			continue
		}
		w := behaviorWeight(b.Type)
		for _, t := range v.Tags {
			counts[t] += w
		}
	}
	type pair struct {
		tag string
		val int
	}
	arr := make([]pair, 0, len(counts))
	for k, v := range counts {
		arr = append(arr, pair{tag: k, val: v})
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].val > arr[j].val })
	if limit <= 0 || limit > len(arr) {
		limit = len(arr)
	}
	out := make([]string, 0, limit)
	for i := 0; i < limit; i++ {
		out = append(out, arr[i].tag)
	}
	return out
}
