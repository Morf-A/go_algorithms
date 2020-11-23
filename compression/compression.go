package compression

import (
	"bufio"
	"bytes"
	"container/list"
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

type bit int8

type TNode struct {
	Element byte
	Name    string
	Count   int
	Left    *TNode
	Right   *TNode
}

type HuffmanCode struct {
	Element byte
	Code    []bit
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

func extend(b []bit, add bit) []bit {
	new := make([]bit, len(b))
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
		code []bit
	}
	var queue []qe
	var ht []HuffmanCode
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
			ht = append(ht, HuffmanCode{
				Element: n.Element,
				Code:    e.code,
			})
		}
	}
	return HuffmanTable(ht)
}

func HuffmanDecode(r io.Reader, tree *TNode) io.Reader {
	return bytes.NewReader(nil)
}

type HuffmanReader struct {
	bitList *list.List
	Buffer  *bufio.Reader
	EOF     bool
}

func (ht *HuffmanReader) PlainToBits(b byte) []bit {
	return nil
}

func (hr *HuffmanReader) nextEncodedBit() (bit, bool) {
	nextPtr := hr.bitList.Front().Value
	if nextPtr == nil {
		plainByte, err := hr.Buffer.ReadByte()
		if err == io.EOF {
			return 0, false
		}
		if err != nil {
			panic(err)
		}
		bits := hr.PlainToBits(plainByte)
		for _, b := range bits {
			hr.bitList.PushBack(b)
		}
		nextPtr = hr.bitList.Front().Value
	}
	return *(nextPtr.(*bit)), true
}

func (hr *HuffmanReader) nextEncodedByte() (byte, bool) {
	var res byte
	for i := 0; i < 8; i++ {
		nextBit, ok := hr.nextEncodedBit()
		if !ok {

		}
		res |= (byte(nextBit) << i)
	}
	return res, true
}

func (hr *HuffmanReader) Read(toFill []byte) (int, error) {
	i := 0
	for i < len(toFill) {
		b, ok := hr.nextEncodedByte()
		if !ok {
			break
		}
		toFill[i] = b
		i++
	}
	var err error
	if i == 0 {
		err = io.EOF
	}
	return i, err
}

func HuffmanEncode(r io.Reader, tree *TNode) (io.Reader, int, int) {
	return &HuffmanReader{
		bitList: list.New(),
		Buffer:  bufio.NewReader(r),
	}, 4324, 420
}

func HuffmanTreeFromTable(ht HuffmanTable) *TNode {
	root := new(TNode)
	for _, hCode := range ht {
		n := root
		for _, b := range hCode.Code {
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
		n.Element = hCode.Element
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
