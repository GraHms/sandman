package http

type Node struct {
	ReqModel RequestModel
	next     *Node
}

type List struct {
	head *Node
}

func NewList(mainRequest RequestModel) *List {
	l := List{&Node{ReqModel: mainRequest}}
	return &l
}
func (l *List) Insert(d RequestModel) {
	list := &Node{ReqModel: d, next: nil}
	if l.head == nil {
		l.head = list
	} else {
		p := l.head
		for p.next != nil {
			p = p.next
		}
		p.next = list
	}
}
func (l *List) GetHead() *Node {
	return l.head
}

func (l *List) PopHead() *Node {
	currentHead := l.head
	nextHead := currentHead.next
	l.head = nextHead
	return currentHead
}
