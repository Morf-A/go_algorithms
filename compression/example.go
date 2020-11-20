package compression

import (
	"fmt"
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

	for k, v := range stat {
		fmt.Println(string(k), v)
	}
	tree := HuffmanTreeFromStat(stat)
	table := HuffmanTreeToTable(tree)
	encodedTable := table.Encode()
	for _, t := range DecodeHuffmanTable(encodedTable) {
		fmt.Println(string(t.Element), t.Code)
	}

}
