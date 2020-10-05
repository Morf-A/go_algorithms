package str

func SlowMaxSequence(a, b string) string {
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
