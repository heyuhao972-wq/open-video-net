package pipeline

import (
	"encoding/json"
	"net/http"

	"recommendation-platform/config"
	"recommendation-platform/model"
	"recommendation-platform/repository"
)

type InNetworkSource struct{}

func (s *InNetworkSource) Name() string { return "in_network" }

func (s *InNetworkSource) GetCandidates(q Query) ([]Candidate, error) {
	// MVP: approximate "in-network" as recently watched items for the user
	follows := q.Follows
	if follows == nil {
		follows = repository.GetFollows()
	}
	if len(follows) == 0 {
		return []Candidate{}, nil
	}

	following := map[string]bool{}
	for _, f := range follows {
		if f.UserID == q.UserID {
			if f.Active {
				following[f.AuthorID] = true
			}
		}
	}
	if len(following) == 0 {
		return []Candidate{}, nil
	}

	videos := LoadIndexVideos(q)
	out := make([]Candidate, 0)
	for _, v := range videos {
		if following[v.AuthorID] {
			out = append(out, Candidate{Video: v})
		}
	}
	return out, nil
}

type IndexSource struct{}

func (s *IndexSource) Name() string { return "index_all" }

func (s *IndexSource) GetCandidates(q Query) ([]Candidate, error) {
	videos := LoadIndexVideos(q)
	out := make([]Candidate, 0, len(videos))
	for _, v := range videos {
		out = append(out, Candidate{Video: v})
	}
	return out, nil
}

type TagSearchSource struct{}

func (s *TagSearchSource) Name() string { return "tag_search" }

func (s *TagSearchSource) GetCandidates(q Query) ([]Candidate, error) {
	tags := q.TopTags
	if len(tags) == 0 {
		return []Candidate{}, nil
	}

	cfg := config.LoadConfig()
	out := make([]Candidate, 0)
	seen := map[string]bool{}

	for _, tag := range tags {
		url := cfg.IndexBase + "/search?q=" + tag
		res, err := http.Get(url)
		if err != nil {
			continue
		}
		if res.StatusCode != http.StatusOK {
			res.Body.Close()
			continue
		}
		var payload struct {
			Videos []model.Video `json:"videos"`
		}
		if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
			res.Body.Close()
			continue
		}
		res.Body.Close()
		for _, v := range payload.Videos {
			if !seen[v.ID] {
				seen[v.ID] = true
				out = append(out, Candidate{Video: v})
			}
		}
	}

	return out, nil
}

type GraphSource struct{}

func (s *GraphSource) Name() string { return "graph_neighbors" }

func (s *GraphSource) GetCandidates(q Query) ([]Candidate, error) {
	graph := q.Graph
	if graph == nil {
		graph = BuildVideoGraph(q.IndexVideos)
	}
	if graph == nil {
		return []Candidate{}, nil
	}

	seeds := seedVideos(q)
	if len(seeds) == 0 {
		return []Candidate{}, nil
	}

	index := q.IndexVideos
	if index == nil {
		index = map[string]model.Video{}
		for _, v := range LoadIndexVideos(q) {
			index[v.ID] = v
		}
	}

	seen := map[string]bool{}
	out := make([]Candidate, 0)
	maxPerSeed := 12

	for _, seed := range seeds {
		neighbors := graph.Neighbors(seed)
		count := 0
		for _, id := range neighbors {
			if seen[id] {
				continue
			}
			v, ok := index[id]
			if !ok {
				continue
			}
			seen[id] = true
			out = append(out, Candidate{Video: v})
			count++
			if count >= maxPerSeed {
				break
			}
		}
	}

	return out, nil
}

func LoadIndexVideos(q Query) []model.Video {
	if q.IndexVideos != nil && len(q.IndexVideos) > 0 {
		out := make([]model.Video, 0, len(q.IndexVideos))
		for _, v := range q.IndexVideos {
			out = append(out, v)
		}
		return out
	}
	cfg := config.LoadConfig()
	url := cfg.IndexBase + "/videos"

	res, err := http.Get(url)
	if err != nil {
		return repository.GetVideos()
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return repository.GetVideos()
	}
	var payload struct {
		Videos []model.Video `json:"videos"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return repository.GetVideos()
	}
	if len(payload.Videos) == 0 {
		return repository.GetVideos()
	}
	return payload.Videos
}
