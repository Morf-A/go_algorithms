package graph

func TopologicalSort(g map[string][]string) []string {
	inDegree := make(map[string]int, len(g))
	var next ListP
	for v := range g {
		inDegree[v] = 0
	}

	for _, adjs := range g {
		for _, v := range adjs {
			inDegree[v]++
		}
	}

	var result []string
	for v := range g {
		if inDegree[v] == 0 {
			next.Push(v)
		}
	}

	for !next.Empty() {
		v := next.Pop()
		for _, aj := range g[v] {
			inDegree[aj]--
			if inDegree[aj] == 0 {
				next.Push(aj)
			}
		}
		result = append(result, v)
	}

	return result
}

func GetShortestDAGPath(s string, g *Graph) *Path {
	p := NewPath(s)
	for _, u := range TopologicalSort(g.GetAdjList()) {
		for _, v := range g.GetAdjList()[u] {
			p.Relax(u, v, g.GetWeight(u, v))
		}
	}
	return p
}

func Dijkstra(s string, g *Graph) *Path {
	h := NewHeap()
	p := NewPath(s)
	for _, u := range g.GetVertices() {
		if u == s {
			h.Insert(u, 0)
		} else {
			h.InsertInf(u)
		}
	}
	adjList := g.GetAdjList()
	for {
		u, ok := h.ExtractMin()
		if !ok {
			break
		}
		for _, v := range adjList[u] {
			if new, changed := p.Relax(u, v, g.GetWeight(u, v)); changed {
				h.DecreaceTo(v, new)
			}
		}
	}
	return p
}
