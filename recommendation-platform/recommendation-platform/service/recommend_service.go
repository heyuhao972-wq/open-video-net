package service

import (
	"recommendation-platform/model"
	"sort"
)

type RecommendService struct {
	registry *AlgorithmRegistry
}

func NewRecommendService() *RecommendService {

	return &RecommendService{
		registry: NewAlgorithmRegistry(),
	}

}

func (s *RecommendService) Recommend(userID string, recType string, limit int) []model.Video {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if recType == "" {
		recType = "default"
	}
	algo := s.registry.Get(recType)
	return algo.Recommend(userID, limit)
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
