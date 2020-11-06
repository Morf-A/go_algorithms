package str

import "fmt"

func ExampleLCSSlow() {
	fmt.Println(GetLCSSlow("CATCGA", "GTACCGTCA"))
}

func ExampleLCS() {
	fmt.Println(GetLCS("CATCGA", "GTACCGTCA"))
}
