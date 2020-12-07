package compression

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

func ExampleZLW() {
	// bookBytes := []byte("ACAAGGTAGGAAAATGCGAAAGCTTAATTGCGGGA")
	bookBytes := []byte("ABCDABC")
	book := bytes.NewReader(bookBytes)

	encoded := LZWEncode(book)

	// res, err := ioutil.ReadAll(encoded)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(res)
	// return

	encodedCopy := &bytes.Buffer{}
	decoded := LZWDecode(io.TeeReader(encoded, encodedCopy))

	original, err := ioutil.ReadAll(decoded)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(original), string(bookBytes))
	if bytes.Equal(original, bookBytes) {
		fmt.Println("Equal")
	} else {
		fmt.Println("Not equal")
	}
	originalLen := len(original)
	encodedLen := len(encodedCopy.Bytes())
	fmt.Printf(
		"%d -> %d %.2f%%\n",
		originalLen,
		encodedLen,
		float64(encodedLen)*float64(100)/float64(originalLen),
	)
}

func ExampleAdaptiveHuffman() {
	bookBytes := []byte("ACAAGGTAGGAAAATGCGAAAGCTTAATTGCGGGA")
	book := bytes.NewReader(bookBytes)

	encoded := HuffmanAdaptiveEncode(book)

	encodedCopy := &bytes.Buffer{}
	decoded, err := HuffmanAdaptiveDecode(io.TeeReader(encoded, encodedCopy))
	if err != nil {
		panic(err)
	}

	original, err := ioutil.ReadAll(decoded)
	if err != nil {
		panic(err)
	}
	if bytes.Equal(original, bookBytes) {
		fmt.Println("Equal")
	} else {
		fmt.Println("Not equal")
	}
	originalLen := len(original)
	encodedLen := len(encodedCopy.Bytes())
	fmt.Printf(
		"%d -> %d %.2f%%\n",
		originalLen,
		encodedLen,
		float64(encodedLen)*float64(100)/float64(originalLen),
	)

}

func ExampleHuffman() {

	book := strings.NewReader("ACAAGGTAGGAAAATGCGAAAGCTTAATTGCGGGA")

	bookBytes, err := ioutil.ReadAll(book)
	if err != nil {
		panic(err)
	}
	if _, err := book.Seek(0, 0); err != nil {
		panic(err)
	}

	stat, err := ByteStatistic(book)
	if err != nil {
		panic(err)
	}

	tree := HuffmanTreeFromStat(stat)
	table := tree.ToTable()
	encodedTable := table.Encode()
	decodedTable := DecodeHuffmanTable(encodedTable)

	newTree := decodedTable.ToTree()

	if _, err := book.Seek(0, 0); err != nil {
		panic(err)
	}
	encoded := HuffmanEncode(book, decodedTable)
	if _, err := book.Seek(0, 0); err != nil {
		panic(err)
	}

	encodedCopy := &bytes.Buffer{}

	reader, err := HuffmanDecode(io.TeeReader(encoded, encodedCopy), newTree)
	if err != nil {
		panic(err)
	}

	original, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	if bytes.Equal(original, bookBytes) {
		fmt.Println("Equal")
	} else {
		fmt.Println("Not equal")
	}
	originalLen := len(original)
	encodedLen := len(encodedCopy.Bytes())
	fmt.Printf(
		"%d -> %d %.2f%%\n",
		originalLen,
		encodedLen,
		float64(encodedLen)*float64(100)/float64(originalLen),
	)

}
