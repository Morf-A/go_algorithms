package main

import (
	"fmt"
	"sort"

	"./array"
	"./graph"
)

func main() {
	exampleInsertionSort()
}

func exampleInsertionSort() {
	input := []int{4, 7, 1, 0, 1, -5, 2, 2, 100, -8}
	fmt.Println(input)
	array.InsertionSort(input)
	fmt.Println(input)
}

func exampleSelectionSort() {
	input := []int{4, 7, 1, 0, 1, -5, 2, 2, 100, -8}
	fmt.Println(input)
	array.SelectionSort(input)
	fmt.Println(input)
}

func exampleBinarySearch() {
	input := []int{1, 4, 6, 8, 9, 13, 16, 19, 30, 35, 48, 50, 51}

	i := sort.Search(len(input), func(i int) bool { return input[i] >= 35 })
	fmt.Printf("Number %d found at index %d\n", 35, i)

	i = array.BinarySearch(input, 35)
	fmt.Printf("Number %d found at index %d\n", 35, i)

	i = array.Search(len(input), func(i int) bool { return input[i] >= 35 })
	fmt.Printf("Number %d found at index %d\n", 35, i)

	i = array.BinarySearch(input, 8)
	fmt.Printf("Number %d found at index %d\n", 8, i)

	i = array.BinarySearch(input, 5)
	fmt.Printf("Number %d found at index %d\n", 5, i)

	i = array.BinarySearchRec(input, 35, 0, len(input)-1)
	fmt.Printf("Number %d found at index %d\n", 35, i)

	i = array.BinarySearchRec(input, 8, 0, len(input)-1)
	fmt.Printf("Number %d found at index %d\n", 8, i)

	i = array.BinarySearchRec(input, 5, 0, len(input)-1)
	fmt.Printf("Number %d found at index %d\n", 5, i)

}

func exampleFloydWarshall() {
	g := graph.NewGraph()
	g.SetWeight("a", "b", 3)
	g.SetWeight("a", "d", 8)
	g.SetWeight("b", "c", 1)
	g.SetWeight("d", "c", -5)
	g.SetWeight("d", "b", 4)
	g.SetWeight("c", "a", 2)

	vertices := g.GetVertices()
	sort.Strings(vertices)
	fmt.Printf(" ")
	for _, u := range vertices {
		fmt.Printf("  %s", u)
	}
	fmt.Println()
	for _, u := range vertices {
		fmt.Printf("%s ", u)
		for _, v := range vertices {
			w, ok := g.Weigth[graph.Edge{u, v}]
			if ok {
				fmt.Printf("% 1d ", w)
			} else {
				fmt.Print(" ∞ ")
			}
		}
		fmt.Println()
	}

	fmt.Println("-----------------")

	paths := graph.FloydWarshall(g)

	fmt.Printf(" ")
	for _, u := range vertices {
		fmt.Printf("  %s", u)
	}
	fmt.Println()
	for _, u := range vertices {
		fmt.Printf("%s ", u)
		for _, v := range vertices {
			w, ok := paths[u].Shortest[v]
			if ok {
				fmt.Printf("% 1d ", w)
			} else {
				fmt.Print(" ∞ ")
			}
		}
		fmt.Println()
	}

	fmt.Println("=================")

	fmt.Printf(" ")
	for _, u := range vertices {
		fmt.Printf("  %s", u)
	}
	fmt.Println()
	for _, u := range vertices {
		fmt.Printf("%s ", u)
		for _, v := range vertices {
			pred, ok := paths[u].Pred[v]
			if ok {
				fmt.Printf(" %s ", pred)
			} else {
				fmt.Print(" - ")
			}
		}
		fmt.Println()
	}

}

func exampleBellmanFord() {
	g := graph.NewGraph()
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
	path := graph.BellmanFord(start, g)
	for _, v := range g.GetVertices() {
		if v != start {
			path.Print(v)
			fmt.Println()
		}
	}
	fmt.Print("----------------------------\n")
	g.SetWeight("t", "s", -4)

	newPath := graph.BellmanFord(start, g)
	cycle := graph.FindNegativeCycle(g, newPath)
	fmt.Println(cycle)

}

func exampleDijkstra() {
	g := graph.NewGraph()
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
	path := graph.Dijkstra(start, g)
	for _, v := range g.GetVertices() {
		if v != start {
			path.Print(v)
			fmt.Println()
		}
	}
}

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
	sorted := graph.TopologicalSort(g)
	fmt.Println(sorted)
}

func exampleHeap() {
	v := graph.NewHeap()
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

func exampleShortestDAGPath() {
	g := graph.NewGraph()
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
	path := graph.GetShortestDAGPath(start, g)
	for _, v := range graph.TopologicalSort(g.GetAdjList()) {
		if v != start {
			path.Print(v)
			fmt.Println()
		}
	}
}
