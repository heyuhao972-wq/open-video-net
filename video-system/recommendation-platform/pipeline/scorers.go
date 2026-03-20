package pipeline

import (
	"time"

	"recommendation-platform/repository"
)

type ViewsScorer struct{}

func (s *ViewsScorer) Name() string { return "views" }

func (s *ViewsScorer) Score(q Query, in []Candidate) []Candidate {
	out := make([]Candidate, 0, len(in))
	for _, c := range in {
		c.Score += c.Video.Views
		out = append(out, c)
	}
	return out
}

type BehaviorScorer struct{}

func (s *BehaviorScorer) Name() string { return "behavior" }

func (s *BehaviorScorer) Score(q Query, in []Candidate) []Candidate {
	boost := map[string]int{}
	behaviors := q.Behaviors
	if behaviors == nil {
		behaviors = repository.GetBehaviors()
	}
	for _, b := range behaviors {
		boost[b.VideoID] += behaviorWeight(b.Type)
	}
	out := make([]Candidate, 0, len(in))
	for _, c := range in {
		c.Score += boost[c.Video.ID]
		out = append(out, c)
	}
	return out
}

type TagSimilarityScorer struct{}

func (s *TagSimilarityScorer) Name() string { return "tag_similarity" }

func (s *TagSimilarityScorer) Score(q Query, in []Candidate) []Candidate {
	if len(q.TopTags) == 0 {
		return in
	}
	top := map[string]bool{}
	for _, t := range q.TopTags {
		top[t] = true
	}
	out := make([]Candidate, 0, len(in))
	for _, c := range in {
		match := 0
		for _, t := range c.Video.Tags {
			if top[t] {
				match++
			}
		}
		c.Score += match * 2
		out = append(out, c)
	}
	return out
}

type RecencyScorer struct {
	MaxDays int
}

func (s *RecencyScorer) Name() string { return "recency" }

func (s *RecencyScorer) Score(q Query, in []Candidate) []Candidate {
	if s.MaxDays <= 0 {
		return in
	}
	out := make([]Candidate, 0, len(in))
	for _, c := range in {
		if c.Video.CreatedAt == 0 {
			out = append(out, c)
			continue
		}
		ageDays := int((time.Now().Unix() - c.Video.CreatedAt) / 86400)
		if ageDays < 0 {
			ageDays = 0
		}
		if ageDays < s.MaxDays {
			c.Score += s.MaxDays - ageDays
		}
		out = append(out, c)
	}
	return out
}

type GraphProximityScorer struct {
	BoostPerEdge int
}

func (s *GraphProximityScorer) Name() string { return "graph_proximity" }

func (s *GraphProximityScorer) Score(q Query, in []Candidate) []Candidate {
	boost := s.BoostPerEdge
	if boost == 0 {
		boost = 3
	}

	graph := q.Graph
	if graph == nil {
		graph = BuildVideoGraph(q.IndexVideos)
	}

	seeds := seedVideos(q)
	if graph == nil || len(seeds) == 0 {
		return in
	}

	neighborBoost := map[string]int{}
	for _, seed := range seeds {
		for _, n := range graph.Neighbors(seed) {
			neighborBoost[n] += boost
		}
	}

	out := make([]Candidate, 0, len(in))
	for _, c := range in {
		c.Score += neighborBoost[c.Video.ID]
		out = append(out, c)
	}
	return out
}

type AuthorDiversityScorer struct {
	PenaltyPerExtra int
}

func (s *AuthorDiversityScorer) Name() string { return "author_diversity" }

func (s *AuthorDiversityScorer) Score(q Query, in []Candidate) []Candidate {
	penalty := s.PenaltyPerExtra
	if penalty <= 0 {
		return in
	}
	counts := map[string]int{}
	out := make([]Candidate, 0, len(in))
	for _, c := range in {
		author := c.Video.AuthorID
		if author == "" {
			author = c.Video.ID
		}
		counts[author]++
		if counts[author] > 1 {
			c.Score -= penalty * (counts[author] - 1)
		}
		out = append(out, c)
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
	case "not_interested":
		return -5
	default:
		return 0
	}
}
