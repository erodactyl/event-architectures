package ringbuffer

import "sync"

type RingBuffer[T any] struct {
	items    []T
	head     int // read from head
	tail     int // write to tail
	size     int
	capacity int
	mu       sync.Mutex
}

func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	return &RingBuffer[T]{
		items:    make([]T, capacity),
		capacity: capacity,
	}
}

func (q *RingBuffer[T]) doubleCapacity() {
	newCapacity := q.capacity * 2
	newItems := make([]T, newCapacity)

	tail := 0
	for i := 0; i < q.size; i++ {
		newItems[i] = q.items[(q.head+i)%q.capacity]
		tail++
	}

	q.items = newItems
	q.capacity = newCapacity
	q.head = 0
	q.tail = tail
}

func (q *RingBuffer[T]) Enqueue(t T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// Check if queue is full
	if q.size+1 == q.capacity {
		q.doubleCapacity()
	}

	q.items[q.tail] = t
	q.tail = (q.tail + 1) % q.capacity
	q.size++
}

func (q *RingBuffer[T]) Dequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	var zero T

	// Check if queue is empty
	if q.head == q.tail {
		return zero, false
	}

	item := q.items[q.head]
	q.items[q.head] = zero // zero out the old item for GC
	q.head = (q.head + 1) % q.capacity
	q.size--
	return item, true
}

func (q *RingBuffer[T]) Size() int {
	return q.size
}
