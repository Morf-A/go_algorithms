package compression

import (
	"bytes"
	"encoding/gob"
)

type TNode struct {
	Element byte
	Count   int
	Left    *TNode
	Right   *TNode
}

type HuffmanTable map[byte][]bit

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

func extend(b []bit, add bit) []bit {
	new := make([]bit, len(b))
	copy(new, b)
	new = append(new, add)
	return new
}

func (n *TNode) ToTable() HuffmanTable {
	if n == nil {
		return nil
	}
	type qe struct {
		node *TNode
		code []bit
	}
	var queue []qe
	ht := make(map[byte][]bit)
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
			ht[n.Element] = e.code
		}
	}
	return HuffmanTable(ht)
}

func (ht HuffmanTable) ToTree() *TNode {
	root := new(TNode)
	for element, bits := range ht {
		n := root
		for _, b := range bits {
			if b == 1 {
				if n.Left == nil {
					n.Left = new(TNode)
				}
				n = n.Left
			} else {
				if n.Right == nil {
					n.Right = new(TNode)
				}
				n = n.Right
			}
		}
		n.Element = element
	}
	return root
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
			Count: a.Count + b.Count,
			Right: a,
			Left:  b,
		})
	}
}
