package lib

import (
	"container/list"
	"sync"

	"github.com/raditya-pratama/go-agent/entity"
)

type Queue struct {
	list *list.List
	sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		list: list.New(),
	}
}

func (q *Queue) Insert(data entity.ActivityLog) {
	q.Lock()
	defer q.Unlock()

	q.list.PushBack(data)
}

func (q *Queue) getFront() *list.Element {
	return q.list.Front()
}

func (q *Queue) GetFront() entity.ActivityLog {
	q.Lock()
	defer q.Unlock()

	data := q.getFront()
	return data.Value.(entity.ActivityLog)
}

func (q *Queue) GetTotal() int {
	q.Lock()
	defer q.Unlock()

	return q.list.Len()
}

func (q *Queue) ReleaseData() {
	q.Lock()
	defer q.Unlock()

	q.list.Remove(q.getFront())
}
