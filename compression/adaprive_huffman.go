package compression

import (
	"bufio"
	"container/list"
	"io"
)

func HuffmanAdaptiveEncode(r io.Reader) io.Reader {
	return &HuffmanAdaptiveEncoder{
		bitList: list.New(),
		state:   hInProgress,
		buffer:  bufio.NewReader(r),
		table:   make(HuffmanTable),
		stat:    make(map[byte]int),
		escape:  []bit{0, 0, 0, 0, 0, 0, 0, 0},
	}
}

type HuffmanAdaptiveEncoder struct {
	bitList  *list.List
	buffer   *bufio.Reader
	table    HuffmanTable
	state    hState
	lastByte byte
	stat     map[byte]int
	escape   []bit
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
	he.stat[b]++
	bits, ok := he.table[b]
	var res []bit
	if !ok {
		res = extend(he.escape, ByteToBits(b)...)
	} else {
		res = bits
	}
	he.table = HuffmanTreeFromStat(he.stat).ToTable()
	return res
}

func HuffmanAdaptiveDecode(r io.Reader) (io.Reader, error) {
	bufReader := bufio.NewReader(r)
	twoBytes := make([]byte, 2)
	_, err := io.ReadFull(bufReader, twoBytes) // we expect at least 2 bytes
	if err != nil {
		return nil, err
	}
	return &HuffmanAdaptiveDecoder{
		bitList:    list.New(),
		buffer:     bufReader,
		tree:       nil,
		node:       nil,
		last2Bytes: [2]byte{twoBytes[0], twoBytes[1]},
	}, nil
}

type HuffmanAdaptiveDecoder struct {
	bitList    *list.List
	buffer     *bufio.Reader
	tree       *TNode
	node       *TNode
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
			break
		}
		toFill[i] = b
		i++
	}
	return
}

func (hd *HuffmanAdaptiveDecoder) nextDecodedByte() (byte, error) {
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
