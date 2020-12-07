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
		if err == io.EOF { //flush maxSeq
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
		maxCode:   code,
		byteList:  list.New(),
		seqLookUp: seqLookUp,
	}
}

type LZWDecoder struct {
	byteList  *list.List
	maxCode   uint32
	seqLookUp map[string]uint32
	sequences []string
	buffer    *bufio.Reader
	maxSeq    []byte
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
	for i := 0; i < 4; i++ {
		nextByte, err := lzwd.buffer.ReadByte()
		if err != nil {
			return 0, err
		}
		intBytes[i] = nextByte
	}
	return binary.BigEndian.Uint32(intBytes), nil
}

func (lzwd *LZWDecoder) nextDecodedBytes() ([]byte, error) {
	nextCode, err := lzwd.nextUint32()
	if err != nil {
		return nil, err
	}
	if int(nextCode) > len(lzwd.sequences)-1 {
		return []byte{0}, nil
	}
	seq := lzwd.sequences[nextCode]
	return []byte(seq), nil
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
	seqLookUp := make(map[string]uint32)
	var sequences []string
	var code uint32 = 0
	for code < 256 {
		seqLookUp[string(code)] = code
		sequences = append(sequences, string(code))
		code++
	}
	return &LZWDecoder{
		buffer:    bufio.NewReader(in),
		maxCode:   code,
		byteList:  list.New(),
		seqLookUp: seqLookUp,
		sequences: sequences,
	}
}
