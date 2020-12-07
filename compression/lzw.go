package compression

import (
	"bufio"
	"container/list"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type LZWEncoder struct {
	bitList   *list.List
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
	next := lzwe.bitList.Front()
	if next == nil {
		bytes, err := lzwe.nextEncodedBytes()
		if err != nil {
			return 0, err
		}
		for _, b := range bytes {
			lzwe.bitList.PushBack(b)
		}
		next = lzwe.bitList.Front()
	}
	res := next.Value.(byte)
	lzwe.bitList.Remove(next)
	return res, nil
}

func (lzwe *LZWEncoder) nextEncodedBytes() ([]byte, error) {
	for {
		if lzwe.isEOF {
			return nil, io.EOF
		}
		nextByte, err := lzwe.buffer.ReadByte()
		if err == io.EOF { //flush maxSeq
			lzwe.isEOF = true
			code, ok := lzwe.seqLookUp[string(lzwe.maxSeq)]
			if !ok {
				return nil, errors.New("Can`t find code by string " + string(lzwe.maxSeq))
			}
			fmt.Println("encode", string(lzwe.maxSeq), code)
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
			fmt.Println("encode", string(lzwe.maxSeq[:len(lzwe.maxSeq)-1]), code)
			fmt.Println("add", string(lzwe.maxSeq), lzwe.maxCode)
			lzwe.maxSeq = []byte{nextByte}
			res := make([]byte, 4)
			binary.BigEndian.PutUint32(res, code)
			return res, nil
		}
	}
}

func initSeq() (uint32, map[string]uint32) {
	seqLookUp := make(map[string]uint32)
	var i uint32 = 0
	for i < 256 {
		seqLookUp[string(i)] = i
		i++
	}
	return i, seqLookUp
}

func LZWEncode(in io.Reader) io.Reader {
	code, seqLookUp := initSeq()
	return &LZWEncoder{
		buffer:    bufio.NewReader(in),
		maxCode:   code,
		bitList:   list.New(),
		seqLookUp: seqLookUp,
	}
}

type LZWDecoder struct {
}

func (lzwd *LZWDecoder) Read(p []byte) (int, error) {
	return 0, nil
}

func LZWDecode(in io.Reader) io.Reader {
	return &LZWDecoder{}
}
