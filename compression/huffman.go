package compression

import (
	"bufio"
	"bytes"
	"container/list"
	"encoding/gob"
	"fmt"
	"io"
	"math"
	"sort"
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

func extend(b []bit, add ...bit) []bit {
	new := make([]bit, len(b))
	copy(new, b)
	new = append(new, add...)
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
			if a != nil && a.Left == nil && a.Right == nil { //if a is root
				return &TNode{Right: a}
			}
			return a
		}

		pq.Insert(&TNode{
			Count: a.Count + b.Count,
			Right: a,
			Left:  b,
		})
	}
}

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
	keys := make([]byte, 0, len(stat))
	for k := range stat {
		keys = append(keys, k)
	}
	//to get the same trees
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	for _, k := range keys {
		pq.Insert(&TNode{
			Count:   stat[k],
			Element: k,
		})
	}
	return &pq
}

type bit int8

func HuffmanDecode(r io.Reader, tree *TNode) (io.Reader, error) {
	bufReader := bufio.NewReader(r)
	twoBytes := make([]byte, 2)
	_, err := io.ReadFull(bufReader, twoBytes) // we expect at least 2 bytes
	if err != nil {
		return nil, err
	}
	return &HuffmanDecoder{
		bitList:    list.New(),
		buffer:     bufReader,
		tree:       tree,
		node:       tree,
		last2Bytes: [2]byte{twoBytes[0], twoBytes[1]},
	}, nil
}

type HuffmanDecoder struct {
	bitList    *list.List
	buffer     *bufio.Reader
	tree       *TNode
	node       *TNode
	isEOF      bool
	last2Bytes [2]byte
}

func ByteToBits(b byte) []bit {
	res := make([]bit, 8)
	for i := 7; i >= 0; i-- {
		if (b & (1 << i)) > 0 {
			res[7-i] = 1
		} else {
			res[7-i] = 0
		}
	}
	return res
}

func (hd *HuffmanDecoder) nextBit() (bit, error) {
	nextPtr := hd.bitList.Front()
	if nextPtr == nil {
		if hd.isEOF {
			return 0, io.EOF
		}
		nextByte, err := hd.buffer.ReadByte()
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF {
			//so, hd.last2Bytes[1] - bits count in hd.last2Bytes[0]
			if hd.last2Bytes[1] == 0 {
				return 0, io.EOF
			}
			bits := ByteToBits(hd.last2Bytes[0])
			for i := 0; i < int(hd.last2Bytes[1]); i++ {
				hd.bitList.PushBack(bits[i])
			}
			hd.isEOF = true
		} else {
			oneByte := hd.last2Bytes[0]
			hd.last2Bytes[0] = hd.last2Bytes[1]
			hd.last2Bytes[1] = nextByte
			for _, oneBit := range ByteToBits(oneByte) {
				hd.bitList.PushBack(oneBit)
			}
		}
		nextPtr = hd.bitList.Front()
	}
	res := nextPtr.Value.(bit)
	hd.bitList.Remove(nextPtr)
	return res, nil
}

func (hd *HuffmanDecoder) nextDecodedByte() (byte, error) {
	for {
		b, err := hd.nextBit()
		if err != nil {
			return 0, err
		}
		if b == 0 {
			hd.node = hd.node.Right
		} else {
			hd.node = hd.node.Left
		}
		if hd.node.Left == nil && hd.node.Right == nil { //leaf
			res := hd.node.Element
			hd.node = hd.tree
			return res, nil
		}
	}
}

func (hd *HuffmanDecoder) Read(toFill []byte) (i int, err error) {
	if hd.isEOF {
		err = io.EOF
		return
	}
	for i < len(toFill) {
		var b byte
		b, err = hd.nextDecodedByte()
		if err != nil {
			if err != io.EOF {
				return 0, err
			}
			break
		}
		toFill[i] = b
		i++
	}
	return
}

type hState int

const (
	hInProgress hState = iota
	hSourceEOF
	hEOF
)

type HuffmanEncoder struct {
	bitList  *list.List
	buffer   *bufio.Reader
	table    HuffmanTable
	state    hState
	lastByte byte
}

func (he *HuffmanEncoder) PlainToBits(b byte) []bit {
	bits, ok := he.table[b]
	if !ok {
		panic(fmt.Sprintf("unknown byte %v\n", b))
	}
	return bits
}

func (he *HuffmanEncoder) nextEncodedBit() (bit, bool) {
	nextPtr := he.bitList.Front()
	if nextPtr == nil {
		plainByte, err := he.buffer.ReadByte()
		if err == io.EOF {
			return 0, false
		}
		if err != nil {
			panic(err)
		}
		bits := he.PlainToBits(plainByte)
		for _, b := range bits {
			he.bitList.PushBack(b)
		}
		nextPtr = he.bitList.Front()
	}
	res := nextPtr.Value.(bit)
	he.bitList.Remove(nextPtr)
	return res, true
}

func (he *HuffmanEncoder) nextEncodedByte() (byte, error) {
	if he.state == hEOF {
		return 0, io.EOF
	}
	if he.state == hSourceEOF {
		he.state = hEOF
		return he.lastByte, nil
	}
	var res byte
	for i := 7; i >= 0; i-- {
		nextBit, ok := he.nextEncodedBit()
		if !ok {
			he.state = hSourceEOF
			he.lastByte = byte(7 - i) //write count of bits in last byte
			return res, nil
		}
		res |= (byte(nextBit) << i)
	}
	return res, nil
}

func (he *HuffmanEncoder) Read(toFill []byte) (i int, err error) {
	for i < len(toFill) {
		var b byte
		b, err = he.nextEncodedByte()
		if err != nil {
			break
		}
		toFill[i] = b
		i++
	}
	if i > 0 {
		err = nil
	}
	return
}

func HuffmanEncode(r io.Reader, table HuffmanTable) io.Reader {
	return &HuffmanEncoder{
		bitList: list.New(),
		state:   hInProgress,
		buffer:  bufio.NewReader(r),
		table:   table,
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
