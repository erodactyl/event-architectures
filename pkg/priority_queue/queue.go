package priorityqueue

import (
	"container/heap"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type Event struct {
	ID             string
	VisibilityTime time.Time
}

func (e *Event) Less(that priority) bool {
	event, ok := that.(*Event)
	if !ok {
		log.Panic("Can't convert event")
	}
	// The one that is newer has LESS priority
	return e.VisibilityTime.After(event.VisibilityTime)
}

type Queue interface {
	Enqueue(e *Event)
	Receive(timeout time.Duration) (*Event, bool)
	Acknowledge(id string) error
	Size() int
}

func NewQueue() Queue {
	return &_queue{
		pq:    priorityQueue[*Event]{},
		index: map[string]*item[*Event]{},
		mu:    sync.RWMutex{},
	}
}

type _queue struct {
	pq    priorityQueue[*Event]
	index map[string]*item[*Event]
	mu    sync.RWMutex
}

func (q *_queue) Enqueue(t *Event) {
	q.mu.Lock()
	defer q.mu.Unlock()
	el := q.pq.Add(t)
	q.index[t.ID] = el
}

func (q *_queue) Receive(timeout time.Duration) (*Event, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	item, exists := q.pq.Peek()
	if !exists {
		return nil, false
	}

	if item.value.VisibilityTime.After(time.Now()) {
		return nil, false
	}

	item.value.VisibilityTime = time.Now().Add(timeout)
	heap.Fix(&q.pq, item.index)
	return item.value, true
}

func (q *_queue) Acknowledge(id string) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	item, exists := q.index[id]
	if !exists {
		return errors.New(fmt.Sprintf("Item with id %s does not exist", id))
	}

	heap.Remove(&q.pq, item.index)
	delete(q.index, id)
	return nil
}

func (q *_queue) Size() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.pq.Len()
}
