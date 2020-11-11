package str

import (
	"bytes"
)

func GetLCSSlow(a, b string) string {
	sA := toSequences2(a)
	sB := toSequences(b)
	var max string
	for _, s1 := range sA {
		for _, s2 := range sB {
			if s1 == s2 && len(s1) > len(max) {
				max = s1
			}
		}
	}
	return max
}

func toSequences2(str string) []string {
	var res []string
	l := len([]byte(str))
	max := 1 << l
	for i := 0; i < max; i++ {
		subs := make([]byte, 0, l)
		for k, s := range []byte(str) {
			if i&(1<<k) != 0 {
				subs = append(subs, s)
			}
		}
		res = append(res, string(subs))
	}
	return res
}

type Store struct {
	Arr []string
}

func (s *Store) add(new string) {
	s.Arr = append(s.Arr, new)
}

func rmCharFromStr(inStr string, i int) string {
	in := []byte(inStr)
	res := make([]byte, 0, len(in))
	for k, s := range in {
		if i != k {
			res = append(res, s)
		}
	}
	return string(res)
}

func doToSequences(in string, i int, max int, s *Store) {
	if i == max {
		s.add(in)
		return
	}
	rmin := rmCharFromStr(in, i)
	doToSequences(in, i+1, max, s)
	doToSequences(rmin, i, max-1, s)
}

func toSequences(str string) []string {
	var s Store
	doToSequences(str, 0, len([]byte(str)), &s)
	return s.Arr
}

func GetLCS(xStr, yStr string) string {
	x := []byte(xStr)
	y := []byte(yStr)
	l := getLCSMap(x, y)
	i := len(x)
	j := len(y)
	var res []byte
	for i > 0 && j > 0 {
		if x[i-1] == y[j-1] {
			res = append(res, x[i-1])
			i--
			j--
			continue
		}
		if l[i][j-1] >= l[i-1][j] {
			j--
		} else {
			i--
		}
	}
	return string(revertBytes(res))
}

func revertBytes(b []byte) []byte {
	i := 0
	j := len(b) - 1
	for i < j {
		b[i], b[j] = b[j], b[i]
		i++
		j--
	}
	return b
}

func getLCSMap(x, y []byte) [][]int {
	lenX := len(x)
	lenY := len(y)
	l := make([][]int, lenX+1)
	for i := 0; i <= lenX; i++ {
		l[i] = make([]int, lenY+1)
	}
	for i := 1; i <= lenX; i++ {
		for j := 1; j <= lenY; j++ {
			if x[i-1] == y[j-1] {
				l[i][j] = l[i-1][j-1] + 1
			} else if l[i-1][j] > l[i][j-1] {
				l[i][j] = l[i-1][j]
			} else {
				l[i][j] = l[i][j-1]
			}
		}
	}
	return l
}

func getNeedleFA(needle []byte) map[byte][]int {
	nextState := make(map[byte][]int)
	for _, n := range needle {
		if _, ok := nextState[n]; ok {
			continue
		}
		nextState[n] = make([]int, len(needle)+1)
		for s := 0; s <= len(needle); s++ {
			readed := make([]byte, s)
			copy(readed, needle[:s])
			readed = append(readed, n)
			maxPref := len(readed)
			if maxPref > len(needle) {
				maxPref = len(needle)
			}
			for i := maxPref; i >= 0; i-- {
				if bytes.Equal(readed[len(readed)-i:], needle[:i]) {
					nextState[n][s] = i
					break
				}
			}
		}
	}
	return nextState
}

//find all substrings with finite automaton
func FindSubstringsFA(needleStr, haystackStr string) []int {
	needle, haystack := []byte(needleStr), []byte(haystackStr)
	nextState := getNeedleFA(needle)
	var (
		res []int
		s   int
	)
	if len(needle) == 0 {
		res = append(res, 0)
	}
	for i, h := range haystack {
		if states, ok := nextState[h]; ok {
			s = states[s]
		} else {
			s = 0
		}
		if s == len(needle) {
			res = append(res, i-len(needle)+1)
		}
	}
	return res
}

func FindSubstrings(needleStr, haystackStr string) []int {
	needle, haystack := []byte(needleStr), []byte(haystackStr)
	var res []int
	for i := 0; i <= len(haystack)-len(needle); i++ {
		j := 0
		for j < len(needle) {
			if needle[j] != haystack[j+i] {
				break
			}
			j++
		}
		if j == len(needle) {
			res = append(res, i)
		}
	}
	return res
}

type Transform struct {
	Op   string
	Args []byte
	Cost int
}

func getTransformMap(x, y []byte) [][]Transform {
	cost := map[string]int{
		"none": 0,
		"ins":  2,
		"del":  2,
		"repl": 1,
		"copy": -1,
	}

	tr := make([][]Transform, len(x)+1)

	tr[0] = make([]Transform, len(y)+1)
	tr[0][0] = Transform{Op: "none"}

	for j := 1; j <= len(y); j++ {
		tr[0][j] = Transform{
			Op:   "ins",
			Cost: tr[0][j-1].Cost + cost["ins"],
			Args: []byte{y[j-1]},
		}
	}

	for i := 1; i <= len(x); i++ {
		tr[i] = make([]Transform, len(y)+1)
		tr[i][0] = Transform{
			Op:   "del",
			Cost: tr[i-1][0].Cost + cost["del"],
			Args: []byte{x[i-1]},
		}
	}

	for i := 1; i <= len(x); i++ {
		for j := 1; j <= len(y); j++ {
			if x[i-1] == y[j-1] {
				tr[i][j] = Transform{
					Op:   "copy",
					Cost: tr[i-1][j-1].Cost + cost["copy"],
					Args: []byte{x[i-1]},
				}
			} else {
				tr[i][j] = Transform{
					Op:   "repl",
					Cost: tr[i-1][j-1].Cost + cost["repl"],
					Args: []byte{x[i-1], y[j-1]},
				}
			}

			if (tr[i-1][j].Cost + cost["del"]) < tr[i][j].Cost {
				tr[i][j] = Transform{
					Op:   "del",
					Cost: tr[i-1][j].Cost + cost["del"],
					Args: []byte{x[i-1]},
				}
			}

			if (tr[i][j-1].Cost + cost["ins"]) < tr[i][j].Cost {
				tr[i][j] = Transform{
					Op:   "ins",
					Cost: tr[i][j-1].Cost + cost["ins"],
					Args: []byte{y[j-1]},
				}
			}
		}
	}
	return tr
}

func GetTransforms(xStr, yStr string) []Transform {
	x, y := []byte(xStr), []byte(yStr)
	tr := getTransformMap(x, y)
	i := len(x)
	j := len(y)
	var res []Transform
	for i > 0 || j > 0 {
		res = append(res, tr[i][j])
		if tr[i][j].Op == "copy" || tr[i][j].Op == "repl" {
			i--
			j--
		} else if tr[i][j].Op == "ins" {
			j--
		} else if tr[i][j].Op == "del" {
			i--
		} else {
			panic("unknown operation: " + tr[i][j].Op)
		}
	}
	i = 0
	j = len(res) - 1
	for i < j {
		res[i], res[j] = res[j], res[i]
		i++
		j--
	}
	return res
}
