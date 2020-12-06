package compression

//TODO: merge huffman and adaptive huffman, optimize table re-building
import (
	"bufio"
	"container/heap"
	"container/list"
	"io"
	"sort"
)

const escape = -1

func HuffmanAdaptiveEncode(r io.Reader) io.Reader {
	stat := make(map[int]int)
	stat[escape] = 0
	return &HuffmanAdaptiveEncoder{
		bitList: list.New(),
		state:   hInProgress,
		buffer:  bufio.NewReader(r),
		table:   make(AdaptiveHuffmanTable),
		stat:    stat,
	}
}

type AdaptiveHuffmanTable map[int][]bit

type HuffmanAdaptiveEncoder struct {
	bitList  *list.List
	buffer   *bufio.Reader
	table    AdaptiveHuffmanTable
	state    hState
	lastByte byte
	stat     map[int]int
}

func (he *HuffmanAdaptiveEncoder) Read(toFill []byte) (i int, err error) {
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

func (he *HuffmanAdaptiveEncoder) nextEncodedByte() (byte, error) {
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

func (he *HuffmanAdaptiveEncoder) nextEncodedBit() (bit, bool) {
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

func (he *HuffmanAdaptiveEncoder) PlainToBits(b byte) []bit {
	he.stat[int(b)]++
	bits, ok := he.table[int(b)]
	var res []bit
	if !ok {
		res = extend(he.table[escape], ByteToBits(b)...)
	} else {
		res = bits
	}
	he.table = AdaptiveHuffmanTreeFromStat(he.stat).ToTable()
	return res
}

func (n *TANode) ToTable() AdaptiveHuffmanTable {
	if n == nil {
		return nil
	}
	type qe struct {
		node *TANode
		code []bit
	}
	var queue []qe
	ht := make(map[int][]bit)
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
	return AdaptiveHuffmanTable(ht)
}

func HuffmanAdaptiveDecode(r io.Reader) (io.Reader, error) {
	bufReader := bufio.NewReader(r)
	twoBytes := make([]byte, 2)
	_, err := io.ReadFull(bufReader, twoBytes) // we expect at least 2 bytes
	if err != nil {
		return nil, err
	}
	stat := make(map[int]int)
	stat[escape] = 0
	return &HuffmanAdaptiveDecoder{
		bitList:    list.New(),
		buffer:     bufReader,
		stat:       stat,
		node:       nil,
		last2Bytes: [2]byte{twoBytes[0], twoBytes[1]},
	}, nil
}

type HuffmanAdaptiveDecoder struct {
	bitList    *list.List
	buffer     *bufio.Reader
	stat       map[int]int
	node       *TANode
	isEOF      bool
	last2Bytes [2]byte
}

func (hd *HuffmanAdaptiveDecoder) Read(toFill []byte) (i int, err error) {
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
		hd.stat[int(b)]++
		hd.node = AdaptiveHuffmanTreeFromStat(hd.stat)
		toFill[i] = b
		i++
	}
	return
}

func (hd *HuffmanAdaptiveDecoder) nextDecodedByte() (byte, error) {
	for {
		if hd.node == nil {
			return hd.nextByte()
		}
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
			if hd.node.Element == escape {
				return hd.nextByte()
			}
			return byte(hd.node.Element), nil
		}
	}
}

func (hd *HuffmanAdaptiveDecoder) nextByte() (byte, error) {
	var res byte
	for i := 7; i >= 0; i-- {
		nextBit, err := hd.nextBit()
		if err != nil {
			return 0, err
		}
		res |= (byte(nextBit) << i)
	}
	return res, nil
}

func (hd *HuffmanAdaptiveDecoder) nextBit() (bit, error) {
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

type TANode struct {
	Element int
	Count   int
	Left    *TANode
	Right   *TANode
}

type AdaptivePriorityQueue struct {
	list []*TANode
}

func (pq *AdaptivePriorityQueue) Len() int {
	return len(pq.list)
}

func (pq *AdaptivePriorityQueue) Less(i, j int) bool {
	return pq.list[i].Count < pq.list[j].Count
}

func (pq *AdaptivePriorityQueue) Swap(i, j int) {
	pq.list[i], pq.list[j] = pq.list[j], pq.list[i]
}

func (pq *AdaptivePriorityQueue) Push(x interface{}) {
	pq.list = append(pq.list, x.(*TANode))
}

func (pq *AdaptivePriorityQueue) Pop() interface{} {
	n := len(pq.list)
	res := pq.list[n-1]
	pq.list = pq.list[:n-1]
	return res
}

func StatToAdaptivePriorityQueue(stat map[int]int) *AdaptivePriorityQueue {
	pq := &AdaptivePriorityQueue{}
	keys := make([]int, 0, len(stat))
	for k := range stat {
		keys = append(keys, k)
	}
	//to get the same trees
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	for _, k := range keys {
		heap.Push(pq, &TANode{
			Count:   stat[k],
			Element: k,
		})
	}
	return pq
}

func TryPop(pq heap.Interface) *TANode {
	if pq.Len() == 0 {
		return nil
	}
	return heap.Pop(pq).(*TANode)
}

func AdaptiveHuffmanTreeFromStat(stat map[int]int) *TANode {
	pq := StatToAdaptivePriorityQueue(stat)
	for {
		a := TryPop(pq)
		b := TryPop(pq)
		if b == nil {
			if a.Left == nil && a.Right == nil { //if a is root
				return &TANode{Right: a}
			}
			return a
		}

		heap.Push(pq, &TANode{
			Count: a.Count + b.Count,
			Right: a,
			Left:  b,
		})
	}
}
