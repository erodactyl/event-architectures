package eventqueue

import (
	"errors"
	ringbuffer "eventarch/pkg/ring_buffer"
	"fmt"
	"sync"
)

type queue[T any] interface {
	Enqueue(t T)
	Dequeue() (T, bool)
	Size() int
}

type Event struct {
	QueueName string
	EventType string
	Body      []byte
}

type EventQueue struct {
	Name  string
	queue queue[*Event]
}

func NewEventQueue(name string) *EventQueue {
	rb := ringbuffer.NewRingBuffer[*Event](10)
	return &EventQueue{Name: name, queue: rb}
}

type EventQueueService struct {
	queues map[string]EventQueue
	mu     sync.RWMutex
}

func NewEventQueueService() *EventQueueService {
	return &EventQueueService{
		queues: map[string]EventQueue{},
		mu:     sync.RWMutex{},
	}
}

type CreateQueueArgs struct {
	Name string
}

func (e *EventQueueService) CreateQueue(args *CreateQueueArgs, reply *bool) error {
	e.mu.RLock()
	if _, exists := e.queues[args.Name]; exists {
		return errors.New(fmt.Sprintf("Queue with name %s already exists", args.Name))
	}
	e.mu.RUnlock()

	e.mu.Lock()
	e.queues[args.Name] = *NewEventQueue(args.Name)
	e.mu.Unlock()

	*reply = true

	return nil
}

type PutArgs = Event

func (e *EventQueueService) Put(args *PutArgs, reply *bool) error {
	e.mu.RLock()
	eq, exists := e.queues[args.QueueName]
	if !exists {
		return errors.New(fmt.Sprintf("Queue with name %s does not exist", args.QueueName))
	}
	e.mu.RUnlock()

	eq.queue.Enqueue(args)

	*reply = true

	return nil
}

type PullArgs struct {
	QueueName string
	Limit     int
}

func (e *EventQueueService) Pull(args *PullArgs, reply *[]*Event) error {
	e.mu.RLock()
	eq, exists := e.queues[args.QueueName]
	if !exists {
		return errors.New(fmt.Sprintf("Queue with name %s does not exist", args.QueueName))
	}
	e.mu.RUnlock()

	*reply = []*Event{}

	for range args.Limit {
		event, exists := eq.queue.Dequeue()
		if !exists {
			break
		}
		*reply = append(*reply, event)
	}

	return nil
}
