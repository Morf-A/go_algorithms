package main

import (
	"fmt"

	"./graph"
)

func excampleTolopogicalSort() {
	g := map[string][]string{
		"x": []string{"a", "b"},
		"a": []string{"c"},
		"b": []string{"c", "e"},
		"c": []string{"f"},
		"f": []string{"g"},
		"d": []string{"e"},
		"e": []string{"f"},
	}
	sorted := topologicalSort(g)
	fmt.Println(sorted)
}

func getShortestDAGPath(s string, g *graph.Graph) *graph.Path {
	p := graph.NewPath(s)
	for _, u := range topologicalSort(g.AdjList) {
		for _, v := range g.AdjList[u] {
			p.Relax(u, v, g.GetWeight(u, v))
		}
	}
	return p
}

func exampleShortestDAGPath() {
	g := graph.NewGraph()
	g.AdjList = map[string][]string{
		"r": []string{"s", "t"},
		"s": []string{"t", "x"},
		"t": []string{"x", "y", "z"},
		"x": []string{"y", "z"},
		"y": []string{"z"},
	}
	g.SetWeight("r", "s", 5)
	g.SetWeight("r", "t", 3)

	g.SetWeight("s", "t", 2)
	g.SetWeight("s", "x", 6)

	g.SetWeight("t", "x", 7)
	g.SetWeight("t", "y", 4)
	g.SetWeight("t", "z", 2)

	g.SetWeight("x", "y", -1)
	g.SetWeight("x", "z", 1)

	g.SetWeight("y", "z", -2)

	start := "s"
	path := getShortestDAGPath(start, g)
	for _, v := range topologicalSort(g.AdjList) {
		if v != start {
			path.Print(v)
			fmt.Println()
		}
	}
}

func main() {
	exampleShortestDAGPath()
}

func topologicalSort(g map[string][]string) []string {
	inDegree := make(map[string]int, len(g))
	var next graph.ListP
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
