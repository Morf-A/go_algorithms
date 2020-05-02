package graph

type Graph struct {
	AdjList map[string][]string
	Weigth  map[Edge]int
}

type Edge [2]string

func NewGraph() *Graph {
	g := &Graph{
		AdjList: make(map[string][]string),
		Weigth:  make(map[Edge]int),
	}
	return g
}

func (g Graph) GetWeight(u, v string) int {
	return g.Weigth[Edge{u, v}]
}

func (g Graph) SetWeight(u, v string, w int) {
	g.Weigth[Edge{u, v}] = w
}
