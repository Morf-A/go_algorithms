package graph

type Graph struct {
	adjList     map[string][]string
	Weight      map[Edge]int
	VerticesSet map[string]bool
}

type Edge [2]string

func NewGraph() *Graph {
	g := &Graph{
		adjList:     make(map[string][]string),
		Weight:      make(map[Edge]int),
		VerticesSet: make(map[string]bool),
	}
	return g
}

func (g Graph) GetWeight(u, v string) int {
	return g.Weight[Edge{u, v}]
}

func (g Graph) GetAdjList() map[string][]string {
	return g.adjList
}

func (g Graph) SetWeight(u, v string, w int) {
	if _, ok := g.Weight[Edge{u, v}]; !ok {
		g.adjList[u] = append(g.adjList[u], v)
		g.VerticesSet[u] = true
		g.VerticesSet[v] = true
	}
	g.Weight[Edge{u, v}] = w
}

func (g Graph) GetVertices() []string {
	var res []string
	for u := range g.VerticesSet {
		res = append(res, u)
	}
	return res
}
