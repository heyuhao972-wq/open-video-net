package pipeline

type Source interface {
	Name() string
	GetCandidates(q Query) ([]Candidate, error)
}

type Filter interface {
	Name() string
	Filter(q Query, in []Candidate) (kept []Candidate, removed []Candidate)
}

type Scorer interface {
	Name() string
	Score(q Query, in []Candidate) []Candidate
}

type Selector interface {
	Name() string
	Select(q Query, in []Candidate, k int) []Candidate
}
