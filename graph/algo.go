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
	adjList := g.GetAdjList()
	for _, u := range TopologicalSort(adjList) {
		for _, v := range adjList[u] {
			p.Relax(u, v, g.GetWeight(u, v))
		}
	}
	return p
}

func FindNegativeCycle(g *Graph, p *Path) []string {
	cyclePathEnd := ""
	for u, adj := range g.GetAdjList() {
		for _, v := range adj {
			if _, changed := p.Relax(u, v, g.GetWeight(u, v)); changed {
				cyclePathEnd = v
			}
		}
	}
	if cyclePathEnd == "" {
		return nil
	}
	v := cyclePathEnd
	visited := make(map[string]bool)
	visited[v] = true
	var cycle []string
	for {
		v = p.Pred[v]
		if visited[v] {
			cycle = append(cycle, v)
			pred := p.Pred[v]
			for v != pred {
				cycle = append(cycle, pred)
				pred = p.Pred[pred]
			}
			break
		}
		visited[v] = true
	}
	revert(cycle)
	return cycle
}

func revert(s []string) []string {
	i := 0
	j := len(s) - 1
	for i < j {
		s[i], s[j] = s[j], s[i]
		i++
		j--
	}
	return s
}

func FloydWarshall(g *Graph) map[string]*Path {
	paths := make(map[string]*Path)

	for u, adj := range g.GetAdjList() {
		paths[u] = NewPath(u)
		for _, v := range adj {
			if w, ok := g.Weigth[Edge{u, v}]; ok {
				paths[u].Shortest[v] = w
				paths[u].Pred[v] = u
			}
		}
	}

	// vertices := g.GetVertices()

	// for n := 1; n < len(vertices); n++ {
	// 	for _, u := range vertices {
	// 		for _, v := range vertices {
	// 			paths[u].Shortest[v]
	// 		}
	// 	}
	// }

	return paths
}

func BellmanFord(s string, g *Graph) *Path {
	p := NewPath(s)

	n := len(g.GetVertices())
	for i := 0; i < n-1; i++ {
		for u, adj := range g.GetAdjList() {
			for _, v := range adj {
				p.Relax(u, v, g.GetWeight(u, v))
			}
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
		vertex, ok := h.ExtractMin()
		if !ok {
			break
		}
		u := vertex.Value
		for _, v := range adjList[u] {
			if new, changed := p.Relax(u, v, g.GetWeight(u, v)); changed {
				h.DecreaseTo(v, new)
			}
		}
	}
	return p
}
