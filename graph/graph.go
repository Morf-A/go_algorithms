package graph

type Graph struct {
	adjList map[string][]string
	Weigth  map[Edge]int
}

type Edge [2]string

func NewGraph() *Graph {
	g := &Graph{
		adjList: make(map[string][]string),
		Weigth:  make(map[Edge]int),
	}
	return g
}

func (g Graph) GetWeight(u, v string) int {
	return g.Weigth[Edge{u, v}]
}

func (g Graph) GetAdjList() map[string][]string {
	return g.adjList
}

func (g Graph) SetWeight(u, v string, w int) {
	g.adjList[u] = append(g.adjList[u], v)
	g.Weigth[Edge{u, v}] = w
}

func (g Graph) GetVertices() []string {
	var res []string
	m := make(map[string]bool)
	for u, list := range g.adjList {
		m[u] = true
		for _, v := range list {
			m[v] = true
		}
	}
	for u := range m {
		res = append(res, u)
	}
	return res
}
