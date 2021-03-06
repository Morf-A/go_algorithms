package graph

import (
	"fmt"
	"log"
	"math"
)

type Vertex struct {
	Key   int
	Value string
}

type Heap struct {
	Array  []Vertex
	Lookup map[string]int
}

func NewHeap() *Heap {
	return &Heap{Lookup: make(map[string]int)}
}

func (hptr *Heap) Check() {
	i := 1
	l := len(hptr.Array)
	for {
		j := 2 * i
		if j > l {
			break
		}
		if hptr.Array[i-1].Key > hptr.Array[j-1].Key {
			log.Println("wrong ", i, j, hptr.Array[i-1], hptr.Array[j-1])
		}
		if j == l {
			break
		}
		j1 := j + 1
		if hptr.Array[i-1].Key > hptr.Array[j1-1].Key {
			log.Println("wrong ", i, j1, hptr.Array[i-1], hptr.Array[j1-1])
		}
		if j1 == l {
			break
		}
		i++
	}
}

func (hptr *Heap) InsertInf(value string) {
	hptr.doInsert(value, math.MaxInt64)
}

func (hptr *Heap) Insert(value string, key int) {
	if key == math.MaxInt64 {
		panic("can`t use max int as key, cause it used as inf value")
	}
	hptr.doInsert(value, key)
}

func (hptr *Heap) doInsert(value string, key int) {
	v := Vertex{key, value}
	hptr.Array = append(hptr.Array, v)
	last := len(hptr.Array) - 1
	hptr.Lookup[v.Value] = last
	hptr.moveUP(last)
}

func (hptr *Heap) DecreaseTo(val string, n int) {
	i := hptr.Lookup[val]
	if hptr.Array[i].Key < n {
		panic(fmt.Sprintf("new value must be less on equal than old value, but %d < %d", hptr.Array[i].Key, n))
	}
	hptr.Array[i].Key = n
	hptr.moveUP(i)
}

func (hptr *Heap) DecreaceKey(val string, n int) {
	i := hptr.Lookup[val]
	hptr.Array[i].Key -= n
	hptr.moveUP(i)
}

func (hptr *Heap) ExtractMin() (Vertex, bool) {
	if len(hptr.Array) == 0 {
		return Vertex{}, false
	}
	res := hptr.Array[0]
	last := len(hptr.Array) - 1
	hptr.swap(last, 0)
	delete(hptr.Lookup, hptr.Array[last].Value)
	hptr.Array = hptr.Array[:last]
	hptr.moveDown(0)
	return res, true
}

func (hptr *Heap) moveUP(i int) {
	h := hptr.Array
	for {
		j := (i - 1) / 2
		if h[i].Key >= h[j].Key {
			break
		}
		hptr.swap(i, j)
		i = j
	}
}

func (hptr *Heap) moveDown(i int) {
	h := hptr.Array
	l := len(h) - 1
	for {
		j := 2*i + 1
		if j > l {
			break
		}
		if j != l && h[j].Key > h[j+1].Key {
			j++
		}
		if h[i].Key <= h[j].Key {
			break
		}
		hptr.swap(i, j)
		i = j
	}
}

func (hptr *Heap) swap(i, j int) {
	h := hptr.Array
	h[i], h[j] = h[j], h[i]
	hptr.Lookup[h[i].Value] = i
	hptr.Lookup[h[j].Value] = j
}
