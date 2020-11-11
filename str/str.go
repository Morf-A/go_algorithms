package str

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
	l := make([][]int, lenX)
	for i := 0; i < lenX; i++ {
		l[i] = make([]int, lenY)
	}
	for i := 1; i < lenX; i++ {
		for j := 1; j < lenY; j++ {
			if x[i] == y[j] {
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
