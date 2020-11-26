package compression

import (
	"fmt"
	"io/ioutil"
	"os"
)

func ExampleHuffman() {
	book, err := os.Open("test.txt")
	if err != nil {
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
	for element, bits := range decodedTable {
		fmt.Println(element, bits)
	}

	newTree := decodedTable.ToTree()
	fmt.Println("--------------")
	for element, bits := range newTree.ToTable() {
		fmt.Println(element, bits)
	}

	if _, err := book.Seek(0, 0); err != nil {
		panic(err)
	}
	encoded := HuffmanEncode(book, decodedTable)

	reader, err := HuffmanDecode(encoded, newTree)
	if err != nil {
		panic(err)
	}

	original, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(original))

}
