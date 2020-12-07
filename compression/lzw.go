package compression

import (
	"bufio"
	"container/list"
	"encoding/binary"
	"errors"
	"io"
)

type LZWEncoder struct {
	byteList  *list.List
	maxCode   uint32
	seqLookUp map[string]uint32
	buffer    *bufio.Reader
	maxSeq    []byte
	isEOF     bool
}

func (lzwe *LZWEncoder) Read(toFill []byte) (i int, err error) {
	for i < len(toFill) {
		var b byte
		b, err = lzwe.nextEncodedByte()
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

func (lzwe *LZWEncoder) nextEncodedByte() (byte, error) {
	next := lzwe.byteList.Front()
	if next == nil {
		bytes, err := lzwe.nextEncodedBytes()
		if err != nil {
			return 0, err
		}
		for _, b := range bytes {
			lzwe.byteList.PushBack(b)
		}
		next = lzwe.byteList.Front()
	}
	res := next.Value.(byte)
	lzwe.byteList.Remove(next)
	return res, nil
}

func (lzwe *LZWEncoder) nextEncodedBytes() ([]byte, error) {
	for {
		if lzwe.isEOF {
			return nil, io.EOF
		}
		nextByte, err := lzwe.buffer.ReadByte()
		if err == io.EOF && len(lzwe.maxSeq) > 0 { //flush maxSeq
			lzwe.isEOF = true
			code, ok := lzwe.seqLookUp[string(lzwe.maxSeq)]
			if !ok {
				return nil, errors.New("Can`t find code by string " + string(lzwe.maxSeq))
			}
			lzwe.maxSeq = nil
			res := make([]byte, 4)
			binary.BigEndian.PutUint32(res, code)
			return res, nil
		}
		if err != nil {
			return nil, err
		}
		lzwe.maxSeq = append(lzwe.maxSeq, nextByte)
		if _, ok := lzwe.seqLookUp[string(lzwe.maxSeq)]; !ok {
			lzwe.maxCode++
			lzwe.seqLookUp[string(lzwe.maxSeq)] = lzwe.maxCode
			code, ok := lzwe.seqLookUp[string(lzwe.maxSeq[:len(lzwe.maxSeq)-1])]
			if !ok {
				return nil, errors.New("Can`t find code by string " + string(lzwe.maxSeq[:len(lzwe.maxSeq)-1]))
			}
			lzwe.maxSeq = []byte{nextByte}
			res := make([]byte, 4)
			binary.BigEndian.PutUint32(res, code)
			return res, nil
		}
	}
}

func LZWEncode(in io.Reader) io.Reader {
	seqLookUp := make(map[string]uint32)
	var code uint32 = 0
	for code < 256 {
		seqLookUp[string(code)] = code
		code++
	}
	return &LZWEncoder{
		buffer:    bufio.NewReader(in),
		maxCode:   code - 1,
		byteList:  list.New(),
		seqLookUp: seqLookUp,
	}
}

type LZWDecoder struct {
	byteList  *list.List
	sequences []string
	buffer    *bufio.Reader
	curSeq    string
	isEOF     bool
}

func (lzwd *LZWDecoder) Read(toFill []byte) (i int, err error) {
	for i < len(toFill) {
		var b byte
		b, err = lzwd.nextDecodedByte()
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

func (lzwd *LZWDecoder) nextUint32() (uint32, error) {
	intBytes := make([]byte, 4)
	_, err := io.ReadFull(lzwd.buffer, intBytes)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(intBytes), nil
}

func (lzwd *LZWDecoder) nextDecodedBytes() ([]byte, error) {
	if lzwd.isEOF {
		return nil, io.EOF
	}
	if lzwd.curSeq == "" {
		curCode, err := lzwd.nextUint32()
		if err != nil {
			return nil, err
		}
		lzwd.curSeq = lzwd.sequences[curCode]
	}

	nextCode, err := lzwd.nextUint32()
	if err == io.EOF {
		lzwd.isEOF = true
		return []byte(lzwd.curSeq), nil
	}
	if err != nil {
		return nil, err
	}
	var nextSeq string
	if int(nextCode) > len(lzwd.sequences)-1 {
		nextSeq = lzwd.curSeq + string(lzwd.curSeq[0])
		lzwd.sequences = append(lzwd.sequences, nextSeq)
	} else {
		nextSeq = lzwd.sequences[nextCode]
		lzwd.sequences = append(lzwd.sequences, lzwd.curSeq+string(nextSeq[0]))
	}
	res := []byte(lzwd.curSeq)
	lzwd.curSeq = nextSeq
	return res, nil
}

func (lzwd *LZWDecoder) nextDecodedByte() (byte, error) {
	next := lzwd.byteList.Front()
	if next == nil {
		bytes, err := lzwd.nextDecodedBytes()
		if err != nil {
			return 0, err
		}
		for _, b := range bytes {
			lzwd.byteList.PushBack(b)
		}
		next = lzwd.byteList.Front()
	}
	res := next.Value.(byte)
	lzwd.byteList.Remove(next)
	return res, nil
}

func LZWDecode(in io.Reader) io.Reader {
	var sequences []string
	for code := uint32(0); code < 256; code++ {
		sequences = append(sequences, string(code))
	}
	return &LZWDecoder{
		buffer:    bufio.NewReader(in),
		byteList:  list.New(),
		sequences: sequences,
	}
}
