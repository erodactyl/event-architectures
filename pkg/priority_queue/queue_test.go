package priorityqueue_test

import (
	priorityqueue "eventarch/pkg/priority_queue"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	q := priorityqueue.NewQueue()

	size := func(s int) {
		if q.Size() != s {
			t.Errorf("Expected size to be %d, got %d", s, q.Size())
		}
	}

	size(0)

	// Earlier, higher priority
	event1 := &priorityqueue.Event{
		ID:             "1",
		VisibilityTime: time.Now().Add(-3 * time.Second),
	}
	event2 := &priorityqueue.Event{
		ID:             "2",
		VisibilityTime: time.Now().Add(-2 * time.Second),
	}

	q.Enqueue(event2)
	q.Enqueue(event1)

	size(2)

	el, exists := q.Receive(2 * time.Second)
	if !exists || el == nil {
		t.Error("Expected event, got nil")
	}

	if el.ID != "1" {
		t.Errorf("Expected event 1, got %s", el.ID)
	}

	el, exists = q.Receive(2 * time.Second)
	if !exists || el == nil {
		t.Error("Expected event, got nil")
	}

	if el.ID != "2" {
		t.Errorf("Expected event 2, got %s", el.ID)
	}

	size(2)

	q.Acknowledge("1")
	size(1)
	q.Acknowledge("2")
	size(0)

	event := &priorityqueue.Event{
		ID:             "1",
		VisibilityTime: time.Now().Add(1 * time.Second),
	}
	q.Enqueue(event)

	size(1)

	el, exists = q.Receive(1 * time.Second)
	if el != nil || exists {
		t.Errorf("Expected no events")
	}

	time.Sleep(1 * time.Second)

	el, exists = q.Receive(1 * time.Second)
	if el == nil || !exists {
		t.Errorf("Expected event")
	}

	// Invisible
	el, exists = q.Receive(1 * time.Second)
	if el != nil || exists {
		t.Errorf("Expected no events")
	}

	// Requeue
	time.Sleep(1 * time.Second)

	el, exists = q.Receive(1 * time.Second)
	if el == nil || !exists {
		t.Errorf("Expected event")
	}
}
