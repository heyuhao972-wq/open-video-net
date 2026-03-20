package pipeline

import "recommendation-platform/model"

type VideoGraph struct {
	adj map[string]map[string]struct{}
}

func BuildVideoGraph(videos map[string]model.Video) *VideoGraph {
	g := &VideoGraph{adj: map[string]map[string]struct{}{}}
	if len(videos) == 0 {
		return g
	}

	byAuthor := map[string][]string{}
	byTag := map[string][]string{}

	for id, v := range videos {
		if v.AuthorID != "" {
			byAuthor[v.AuthorID] = append(byAuthor[v.AuthorID], id)
		}
		for _, t := range v.Tags {
			if t == "" {
				continue
			}
			byTag[t] = append(byTag[t], id)
		}
	}

	for _, ids := range byAuthor {
		linkAll(g, ids)
	}
	for _, ids := range byTag {
		linkAll(g, ids)
	}

	return g
}

func linkAll(g *VideoGraph, ids []string) {
	for i := 0; i < len(ids); i++ {
		for j := i + 1; j < len(ids); j++ {
			g.addEdge(ids[i], ids[j])
			g.addEdge(ids[j], ids[i])
		}
	}
}

func (g *VideoGraph) addEdge(from string, to string) {
	if from == "" || to == "" || from == to {
		return
	}
	if g.adj[from] == nil {
		g.adj[from] = map[string]struct{}{}
	}
	g.adj[from][to] = struct{}{}
}

func (g *VideoGraph) Neighbors(id string) []string {
	if g == nil || g.adj == nil {
		return nil
	}
	m := g.adj[id]
	if len(m) == 0 {
		return nil
	}
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func seedVideos(q Query) []string {
	seen := map[string]bool{}
	out := make([]string, 0)

	for _, b := range q.Behaviors {
		if b.VideoID == "" {
			continue
		}
		if !seen[b.VideoID] {
			seen[b.VideoID] = true
			out = append(out, b.VideoID)
		}
	}

	if len(out) == 0 {
		for _, f := range q.Follows {
			if f.AuthorID == "" {
				continue
			}
			for id, v := range q.IndexVideos {
				if v.AuthorID == f.AuthorID && !seen[id] {
					seen[id] = true
					out = append(out, id)
				}
			}
		}
	}

	if len(out) == 0 && len(q.TopTags) > 0 {
		for id, v := range q.IndexVideos {
			if hasTag(v.Tags, q.TopTags) && !seen[id] {
				seen[id] = true
				out = append(out, id)
			}
		}
	}

	return out
}

func hasTag(tags []string, targets []string) bool {
	if len(tags) == 0 || len(targets) == 0 {
		return false
	}
	set := map[string]struct{}{}
	for _, t := range tags {
		set[t] = struct{}{}
	}
	for _, t := range targets {
		if _, ok := set[t]; ok {
			return true
		}
	}
	return false
}
