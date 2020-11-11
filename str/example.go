package str

import "fmt"

func ExampleLCSSlow() {
	fmt.Println(GetLCSSlow("CATCGA", "GTACCGTCA"))
}

func ExampleLCS() {
	fmt.Println(GetLCS("CATCGA", "GTACCGTCA"))
}

func ExampleFindSubstringsFA() {
	fmt.Println(FindSubstringsFA("AAC", "FISAIDSAACACIOAEJAAAC"))
	fmt.Println(FindSubstringsFA("", "ABC"))
	fmt.Println(FindSubstringsFA("ABC", ""))
	fmt.Println(FindSubstringsFA("", ""))
}

func ExampleFindSubstrings() {
	fmt.Println(FindSubstrings("AAC", "FISAIDSAACACIOAEJAAAC"))
	fmt.Println(FindSubstrings("", "ABC"))
	fmt.Println(FindSubstrings("ABC", ""))
	fmt.Println(FindSubstrings("", ""))
}
