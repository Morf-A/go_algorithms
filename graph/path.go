package graph

import "fmt"

type Path struct {
	Vertex   string
	Shortest map[string]int
	Pred     map[string]string
}

func NewPath(s string) *Path {
	p := &Path{
		Vertex:   s,
		Shortest: make(map[string]int),
		Pred:     make(map[string]string),
	}
	p.Shortest[s] = 0
	return p
}

func (p *Path) Relax(u, v string, weightUV int) (int, bool) {
	shortestU, ok := p.Shortest[u]
	if !ok {
		return 0, false
	}
	shortestV, ok := p.Shortest[v]
	if !ok || shortestU+weightUV < shortestV {
		shortestV = shortestU + weightUV
		p.Pred[v] = u
		p.Shortest[v] = shortestV
		return shortestV, true
	}
	return 0, false
}

func (p *Path) Print(to string) {
	fmt.Print("to ", to, ": ")
	var predList []string
	for {
		next, ok := p.Pred[to]
		if !ok {
			break
		}
		predList = append(predList, to)
		to = next
	}

	if len(predList) == 0 {
		fmt.Print("NULL")
	}

	for i := len(predList) - 1; i >= 0; i-- {
		v := predList[i]
		fmt.Printf("[%s:%d] ", v, p.Shortest[v])
	}
}
