package compression

import (
	"fmt"
	"os"
)

func ExampleHuffman() {
	book, err := os.Open("book.txt")
	if err != nil {
		panic(err)
	}
	stat, err := ByteStatistic(book)
	if err != nil {
		panic(err)
	}
	for _, s := range StatToSortedNodes(stat) {
		fmt.Println(string(s.Element), s.Count)
	}

}
