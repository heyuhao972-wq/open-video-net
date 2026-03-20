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

func (s *RecommendService) Recommend(userID string) []model.Video {

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
		K:        20,
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
