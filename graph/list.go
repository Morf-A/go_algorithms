package graph

type ListP struct {
	l *List
}

type List struct {
	e string
	n *List
}

func (p *ListP) Empty() bool {
	return p.l == nil
}

func (p *ListP) Push(e string) {
	p.l = &List{e, p.l}
}

func (p *ListP) Pop() string {
	res := p.l.e
	p.l = p.l.n
	return res
}
