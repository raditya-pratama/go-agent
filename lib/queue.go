package lib

import "fmt"

type Node struct {
	prev *Node
	next *Node
	key  interface{}
}

type List struct {
	head *Node
	tail *Node
}

func NewQueue() *List {
	return &List{}
}

func (L *List) Insert(key interface{}) {
	list := &Node{
		next: L.head,
		key:  key,
	}
	if L.head != nil {
		L.head.prev = list
	}
	L.head = list

	l := L.head
	for l.next != nil {
		l = l.next
	}
	L.tail = l
}

func (l *List) Display() {
	if l.head == nil {
		fmt.Println("list is empty")
		return
	}
	list := l.head
	for list != nil {
		fmt.Printf("%+v ->", list.key)
		list = list.next
		l.head = nil
	}
}

func (l *List) GetHead() *Node {
	return l.head
}
func (l *List) SetHead(node *Node) {
	l.head = node
}

func GetNext(node *Node) *Node {
	return node.next
}

func GetValue(node *Node) interface{} {
	return node.key
}

func Display(list *Node) {
	for list != nil {
		fmt.Printf("%v ->", list.key)
		list = list.next
	}
}

func ShowBackwards(list *Node) {
	for list != nil {
		fmt.Printf("%v <-", list.key)
		list = list.prev
	}
}

func (l *List) Reverse() {
	curr := l.head
	var prev *Node
	l.tail = l.head

	for curr != nil {
		next := curr.next
		curr.next = prev
		prev = curr
		curr = next
	}
	l.head = prev
}
