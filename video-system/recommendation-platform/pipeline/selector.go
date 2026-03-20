package pipeline

import "sort"

type TopKSelector struct{}

func (s *TopKSelector) Name() string { return "top_k" }

func (s *TopKSelector) Select(q Query, in []Candidate, k int) []Candidate {
	if k <= 0 || k > len(in) {
		k = len(in)
	}
	sort.Slice(in, func(i, j int) bool {
		return in[i].Score > in[j].Score
	})
	return in[:k]
}
