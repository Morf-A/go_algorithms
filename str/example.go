package str

import "fmt"

func ExampleLCSSlow() {
	max := GetLCS("CATCGA", "GTACCGTCA")
	fmt.Println(max)
}
