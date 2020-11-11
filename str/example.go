package str

import "fmt"

func ExampleLCSSlow() {
	fmt.Println(GetLCSSlow("CATCGA", "GTACCGTCA"))
}

func ExampleLCS() {
	fmt.Println(GetLCS("CATCGC", "GTACCGTCA"))
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

func ExampleGetTransforms() {
	trs := GetTransforms("ACAAGC", "CCGT")
	for _, tr := range trs {
		args := []string{}
		for _, a := range tr.Args {
			args = append(args, string(a))
		}
		fmt.Print(tr.Op, args, "(", tr.Cost, ")\n")
	}
}
