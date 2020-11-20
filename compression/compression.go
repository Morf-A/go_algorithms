package compression

import (
	"bufio"
	"io"
	"sort"
)

func ByteStatistic(r io.Reader) (map[byte]int, error) {
	m := make(map[byte]int)
	br := bufio.NewReader(r)
	for {
		b, err := br.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		m[b]++
	}
	return m, nil
}

func StatToSortedNodes(stat map[byte]int) []TNode {
	nodes := make([]TNode, len(stat))
	i := 0
	for k, v := range stat {
		nodes[i] = TNode{Element: k, Count: v}
		i++
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Count < nodes[j].Count
	})
	return nodes
}

func extractMin(treeNodesP *[]TNode, statNodesP *[]TNode) {
	treeNodes := *(treeNodesP)
	statNodes := *(statNodesP)
	minTree := (*treeNodes)[0]
	minStat := *statNodes[len(statNodes)-1]
	if minTree < minStat {
		return minTree
	}
}

func CreateHuffmanTree(stat map[byte]int) {
	statNodes := StatToSortedNodes(stat)
	treeNodes := make([]BTNode, 0, 2*len(statNodes))
	t := BinaryTree{}
	for {

	}

}

type StatNode struct {
	Element byte
	Count   int
}
