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
	table := HuffmanTreeToTable(tree)
	encodedTable := table.Encode()
	decodedTable := DecodeHuffmanTable(encodedTable)
	for _, t := range decodedTable {
		fmt.Println(t.Element, t.Code)
	}

	newTree := HuffmanTreeFromTable(decodedTable)
	fmt.Println("--------------")
	for _, t := range HuffmanTreeToTable(newTree) {
		fmt.Println(t.Element, t.Code)
	}

	if _, err := book.Seek(0, 0); err != nil {
		panic(err)
	}
	encoded, old, new := HuffmanEncode(book, newTree)
	fmt.Printf("old: %d, new: %d, efficiency: %.2f%%", old, new, float64(new)/float64(old)*100)

	reader := HuffmanDecode(encoded, newTree)

	original, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(original))

}
