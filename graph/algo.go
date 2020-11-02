package graph

import (
	"fmt"
	"sort"
	"strconv"
)

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

func FindNegativeCycle(g *Graph, start string) []string {
	p := BellmanFord(start, g)
	//find cycle tail (can be outside)
	var cycleTail string
	for u, adj := range g.GetAdjList() {
		for _, v := range adj {
			if _, changed := p.Relax(u, v, g.GetWeight(u, v)); changed {
				cycleTail = v
				break
			}
		}
	}
	if cycleTail == "" {
		return nil
	}
	//find cycle end
	visited := make(map[string]bool)
	var cycleEnd string
	v := cycleTail
	for {
		if visited[v] {
			cycleEnd = v
			break
		}
		visited[v] = true
		v = p.Pred[v]
	}
	//find cycle
	var cycle []string
	v = cycleEnd
	for {
		cycle = append(cycle, v)
		v = p.Pred[v]
		if v == cycleEnd {
			break
		}
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

func IntLeghth(i int) int {
	return len([]byte(strconv.Itoa(i)))
}

func PrintAdjMap(g *Graph) {
	maxLen := 0
	for _, w := range g.Weigth {
		wLength := IntLeghth(w)
		if maxLen < IntLeghth(w) {
			maxLen = wLength
		}
	}
	format := "% " + strconv.Itoa(maxLen+1) + "s" // for example, "% 5s"
	vertices := g.GetVertices()
	sort.Strings(vertices)
	fmt.Printf(" ")
	for _, u := range vertices {
		fmt.Printf(format, u)
	}
	fmt.Println()
	for _, u := range vertices {
		fmt.Printf("%s", u)
		for _, v := range vertices {
			w, ok := g.Weigth[Edge{u, v}]
			if ok {
				fmt.Printf(format, strconv.Itoa(w))
			} else {
				fmt.Printf(format, "∞")
			}
		}
		fmt.Println()
	}
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
