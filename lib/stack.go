package lib

import (
	"fmt"
	"sync"

	"github.com/raditya-pratama/go-agent/entity"
)

type Node struct {
	prev *Node
	next *Node
	data entity.ActivityLog
}

type List struct {
	head *Node
	tail *Node
	sync.Mutex
}

func NewStack() *List {
	return &List{}
}

func (L *List) Insert(key entity.ActivityLog) {
	L.Lock()
	defer L.Unlock()
	list := &Node{
		next: L.head,
		data: key,
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
	l.Lock()
	defer l.Unlock()
	if l.head == nil {
		fmt.Println("list is empty")
		return
	}
	list := l.head
	for list != nil {
		fmt.Printf("%+v ->", list.data)
		list = list.next
		l.head = nil
	}
}

func (l *List) GetHead() *Node {
	l.Lock()
	defer l.Unlock()
	return l.head
}

func (l *List) SetHead(node *Node) {
	l.Lock()
	defer l.Unlock()
	l.head = node
}

func GetNext(node *Node) *Node {

	return node.next
}

func GetValue(node *Node) entity.ActivityLog {
	return node.data
}

func Display(list *Node) {
	for list != nil {
		fmt.Printf("%v ->", list.data)
		list = list.next
	}
}

func ShowBackwards(list *Node) {
	for list != nil {
		fmt.Printf("%v <-", list.data)
		list = list.prev
	}
}

func (l *List) Reverse() {
	l.Lock()
	defer l.Unlock()
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
