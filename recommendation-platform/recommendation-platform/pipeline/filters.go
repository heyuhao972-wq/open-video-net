package pipeline

import (
	"time"

	"recommendation-platform/repository"
)

type PreviouslySeenFilter struct{}

func (f *PreviouslySeenFilter) Name() string { return "previously_seen" }

func (f *PreviouslySeenFilter) Filter(q Query, in []Candidate) ([]Candidate, []Candidate) {
	behaviors := repository.GetBehaviors()
	seen := map[string]bool{}
	for _, b := range behaviors {
		if b.UserID == q.UserID {
			seen[b.VideoID] = true
		}
	}
	if len(seen) == 0 {
		return in, nil
	}

	var kept []Candidate
	var removed []Candidate
	for _, c := range in {
		if seen[c.Video.ID] {
			removed = append(removed, c)
		} else {
			kept = append(kept, c)
		}
	}
	return kept, removed
}

type DedupFilter struct{}

func (f *DedupFilter) Name() string { return "dedup" }

func (f *DedupFilter) Filter(q Query, in []Candidate) ([]Candidate, []Candidate) {
	seen := map[string]bool{}
	var kept []Candidate
	var removed []Candidate
	for _, c := range in {
		if seen[c.Video.ID] {
			removed = append(removed, c)
		} else {
			seen[c.Video.ID] = true
			kept = append(kept, c)
		}
	}
	return kept, removed
}

type AgeFilter struct {
	MaxAgeDays int
}

func (f *AgeFilter) Name() string { return "age" }

func (f *AgeFilter) Filter(q Query, in []Candidate) ([]Candidate, []Candidate) {
	if f.MaxAgeDays <= 0 {
		return in, nil
	}
	cutoff := time.Now().Add(-time.Duration(f.MaxAgeDays) * 24 * time.Hour).Unix()
	var kept []Candidate
	var removed []Candidate
	for _, c := range in {
		if c.Video.CreatedAt > 0 && c.Video.CreatedAt < cutoff {
			removed = append(removed, c)
		} else {
			kept = append(kept, c)
		}
	}
	return kept, removed
}
