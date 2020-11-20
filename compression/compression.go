package compression

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"io"
	"math"
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

func StatToPriorityQueue(stat map[byte]int) *PriorityQueue {
	pq := PriorityQueue{}
	for k, v := range stat {
		pq.Insert(&TNode{
			Count:   v,
			Element: k,
			Name:    string(k),
		})
	}
	return &pq
}

type TNode struct {
	Element byte
	Name    string
	Count   int
	Left    *TNode
	Right   *TNode
}

type HuffmanCode struct {
	Element byte
	Code    []byte
}

type HuffmanTable []HuffmanCode

func (ht HuffmanTable) Encode() []byte {
	buf := bytes.Buffer{}
	if err := gob.NewEncoder(&buf).Encode(ht); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func DecodeHuffmanTable(b []byte) HuffmanTable {
	var ht HuffmanTable
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&ht); err != nil {
		panic(err)
	}
	return ht
}

func extend(b []byte, add byte) []byte {
	new := make([]byte, len(b))
	copy(new, b)
	new = append(new, add)
	return new
}

func HuffmanTreeToTable(n *TNode) HuffmanTable {
	if n == nil {
		return nil
	}
	type qe struct {
		node *TNode
		code []byte
	}
	var queue []qe
	var hCodes []HuffmanCode
	queue = append(queue, qe{node: n, code: nil})
	for len(queue) != 0 {
		e := queue[0]
		n := e.node
		queue = queue[1:]
		if n.Left != nil {
			queue = append(queue, qe{node: n.Left, code: extend(e.code, 1)})
		}
		if n.Right != nil {
			queue = append(queue, qe{node: n.Right, code: extend(e.code, 0)})
		}
		if n.Left == nil && n.Right == nil { //leaf
			hCodes = append(hCodes, HuffmanCode{
				Element: n.Element,
				Code:    e.code,
			})
		}
	}
	return HuffmanTable(hCodes)
}

func HuffmanTreeFromTable(ht HuffmanTable) *TNode {
	return nil
}

func HuffmanTreeFromStat(stat map[byte]int) *TNode {
	pq := StatToPriorityQueue(stat)
	for {
		a := pq.ExtractMin()
		b := pq.ExtractMin()
		if b == nil {
			return a
		}

		pq.Insert(&TNode{
			Name:  a.Name + b.Name,
			Count: a.Count + b.Count,
			Right: a,
			Left:  b,
		})
	}
}

type PriorityQueue struct {
	list []*TNode
}

func (pq *PriorityQueue) Insert(n *TNode) {
	pq.list = append(pq.list, n)
}

func (pq *PriorityQueue) ExtractMin() *TNode {
	if len(pq.list) == 0 {
		return nil
	}
	min := math.MaxInt64
	var minNodeID int
	for i, n := range pq.list {
		if n.Count <= min {
			min = n.Count
			minNodeID = i
		}
	}
	last := len(pq.list) - 1
	pq.list[last], pq.list[minNodeID] = pq.list[minNodeID], pq.list[last]
	res := pq.list[last]
	pq.list = pq.list[:last]
	return res
}
