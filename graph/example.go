package graph

import (
	"fmt"
	"strconv"
)

func ExcampleTolopogicalSort() {
	g := map[string][]string{
		"x": []string{"a", "b"},
		"a": []string{"c"},
		"b": []string{"c", "e"},
		"c": []string{"f"},
		"f": []string{"g"},
		"d": []string{"e"},
		"e": []string{"f"},
	}
	sorted := TopologicalSort(g)
	fmt.Println(sorted)
}

func ExampleFloydWarshall() {
	g := NewGraph()
	g.SetWeight("a", "b", 3)
	g.SetWeight("a", "d", 8)
	g.SetWeight("b", "c", 1)
	g.SetWeight("d", "c", -5)
	g.SetWeight("d", "b", 4)
	g.SetWeight("c", "a", 2)

	fmt.Print("-------Adjacency---------\n\n")
	PrintMap(g.GetVertices(), func(u, v string) string {
		w, ok := g.Weight[Edge{u, v}]
		if !ok {
			return "∞"
		}
		return strconv.Itoa(w)
	})

	paths := FloydWarshall(g)
	fmt.Print("\n-------Shortest----------\n\n")
	PrintMap(g.GetVertices(), func(u, v string) string {
		s, ok := paths[u].Shortest[v]
		if !ok {
			return "∞"
		}
		return strconv.Itoa(s)
	})
	fmt.Print("\n-------Previous----------\n\n")

	PrintMap(g.GetVertices(), func(u, v string) string {
		p, ok := paths[u].Pred[v]
		if !ok {
			return " -"
		}
		return " " + p
	})
}

func ExampleFindNegativeCycle() {
	g := NewGraph()
	g.SetWeight("e", "h", 1)
	g.SetWeight("a", "b", 1)
	g.SetWeight("b", "c", 1)
	g.SetWeight("c", "d", 1)
	g.SetWeight("d", "e", -2)
	g.SetWeight("e", "b", -2)
	g.SetWeight("d", "f", 1)
	g.SetWeight("f", "g", 1)
	cycle := FindNegativeCycle(g, "a")
	fmt.Println(cycle)
}

func ExampleBellmanFord() {
	g := NewGraph()
	g.SetWeight("s", "t", 6)
	g.SetWeight("s", "y", 7)
	g.SetWeight("t", "x", 5)
	g.SetWeight("t", "y", 8)
	g.SetWeight("t", "z", -4)
	g.SetWeight("x", "t", -2)
	g.SetWeight("z", "x", 7)
	g.SetWeight("z", "s", 2)
	g.SetWeight("y", "x", -3)
	g.SetWeight("y", "z", 9)

	start := "s"
	path := BellmanFord(start, g)
	for _, v := range g.GetVertices() {
		if v != start {
			path.Print(v)
			fmt.Println()
		}
	}
}

func ExampleDijkstra() {
	g := NewGraph()
	g.SetWeight("s", "t", 6)
	g.SetWeight("s", "y", 4)
	g.SetWeight("t", "x", 3)
	g.SetWeight("t", "y", 2)
	g.SetWeight("x", "z", 4)
	g.SetWeight("z", "x", 5)
	g.SetWeight("z", "s", 7)
	g.SetWeight("y", "t", 1)
	g.SetWeight("y", "z", 3)
	g.SetWeight("y", "x", 9)

	start := "s"
	path := Dijkstra(start, g)
	for _, v := range g.GetVertices() {
		if v != start {
			path.Print(v)
			fmt.Println()
		}
	}
}

func ExampleHeap() {
	v := NewHeap()
	v.Insert("a", 4)
	v.Insert("b", 10)
	v.Insert("c", 2)
	v.Insert("d", 15)
	v.Insert("e", 8)
	v.Insert("f", 11)
	v.Insert("z", 111)
	v.InsertInf("z")
	v.Insert("g", 16)
	v.Check()
	v.Insert("h", 14)
	v.Insert("i", 18)
	v.InsertInf("t")
	v.Insert("j", 17)
	fmt.Println(v.Array)
	v.ExtractMin()
	fmt.Println(v.Array)
	v.Insert("k", 1)
	v.ExtractMin()
	v.Check()
	v.DecreaceKey("l", 11)
	v.Check()
	v.ExtractMin()
	v.Check()
	fmt.Println(v.Array)
	v.DecreaceKey("f", 3)
	v.Check()
	fmt.Println(v.Array)
	v.ExtractMin()
	fmt.Println(v.Array)
	v.Check()
	v.ExtractMin()
	fmt.Println(v.Array)
}

func ExampleShortestDAGPath() {
	g := NewGraph()
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
	path := GetShortestDAGPath(start, g)
	for _, v := range TopologicalSort(g.GetAdjList()) {
		if v != start {
			path.Print(v)
			fmt.Println()
		}
	}
}
