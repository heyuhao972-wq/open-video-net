package pipeline

import "sort"

type Engine struct {
	Sources  []Source
	Filters  []Filter
	Scorers  []Scorer
	Selector Selector
	K        int
}

type Result struct {
	Retrieved []Candidate
	Filtered  []Candidate
	Selected  []Candidate
}

func (e *Engine) Run(q Query) Result {
	var retrieved []Candidate
	for _, s := range e.Sources {
		cands, err := s.GetCandidates(q)
		if err != nil {
			continue
		}
		retrieved = append(retrieved, cands...)
	}

	kept := retrieved
	var removed []Candidate
	for _, f := range e.Filters {
		k, r := f.Filter(q, kept)
		kept = k
		removed = append(removed, r...)
	}

	scored := kept
	for _, s := range e.Scorers {
		scored = s.Score(q, scored)
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	k := e.K
	if k <= 0 || k > len(scored) {
		k = len(scored)
	}
	selected := scored[:k]

	if e.Selector != nil {
		selected = e.Selector.Select(q, selected, k)
	}

	return Result{
		Retrieved: retrieved,
		Filtered:  removed,
		Selected:  selected,
	}
}
