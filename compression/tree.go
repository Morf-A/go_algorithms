package compression

type Tree struct {
	Root    *TNode
	Current *TNode
}

type TNode struct {
	Element byte
	Count   int
	Left    *TNode
	LPrefix byte
	Right   *TNode
	RPrefix byte
}

func (t *Tree) NewLeft() {

}

func (t *Tree) NewRight() {

}

func (t *Tree) SetValue([]byte) {

}

func (t *Tree) GoRoot() {
	bt.Current = bt.Root
}
