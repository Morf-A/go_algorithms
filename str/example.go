package str

import "fmt"

func ExampleMaxSequence() {
	max := SlowMaxSequence("CATCGA", "GTACCGTCA")
	fmt.Println(max)
}
