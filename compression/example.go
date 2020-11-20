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
	for _, t := range DecodeHuffmanTable(encodedTable) {
		fmt.Println(t.Element, t.Code)
	}

	if _, err := book.Seek(0, 0); err != nil {
		panic(err)
	}

	encoded, old, new := HuffmanEncode(book, tree)
	fmt.Printf("old: %d, new: %d, efficiency: %.2f%%", old, new, float64(new)/float64(old)*100)

	reader := HuffmanDecode(encoded, tree)

	original, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(original))

}
