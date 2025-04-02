package priorityqueue

type priority interface {
	Less(that priority) bool
}

type item[T priority] struct {
	value T
	index int
}

type priorityQueue[T priority] []*item[T]

func newpriorityQueue[T priority]() *priorityQueue[T] {
	pq := &priorityQueue[T]{}
	return pq
}

func (q priorityQueue[T]) Len() int {
	return len(q)
}

func (q priorityQueue[T]) Less(i, j int) bool {
	return q[i].value.Less(q[j].value)
}

func (q priorityQueue[T]) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (pq *priorityQueue[T]) Push(x any) {
	el := x.(*item[T])
	el.index = pq.Len()
	*pq = append(*pq, el)
}

func (pq *priorityQueue[T]) Pop() any {
	n := pq.Len()
	if n == 0 {
		return nil
	}
	last := (*pq)[n-1]
	last.index = -1
	(*pq)[n-1] = nil // clear the reference so GC can reclaim it
	*pq = (*pq)[:n-1]
	return last
}

func (pq *priorityQueue[T]) Add(t T) *item[T] {
	el := &item[T]{value: t}
	pq.Push(el)
	return el
}

func (pq *priorityQueue[T]) Get() (*item[T], bool) {
	var zero *item[T]
	el := pq.Pop()
	if el == nil {
		return zero, false
	}
	return el.(*item[T]), true
}

func (pq *priorityQueue[T]) Peek() (*item[T], bool) {
	n := pq.Len()
	if n == 0 {
		return nil, false
	}
	last := (*pq)[n-1]
	return last, true
}
