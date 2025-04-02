package ringbuffer_test

import (
	ringbuffer "eventbus/pkg/ring_buffer"
	"sync"
	"testing"
)

func TestRingBuffer(t *testing.T) {
	rb := ringbuffer.NewRingBuffer[int](10)
	for i := range 10 {
		rb.Enqueue(i)
	}

	size := rb.Size()
	if size != 10 {
		t.Errorf("Expected size %d, got %d", 10, size)
	}

	for i := range 10 {
		item, exists := rb.Dequeue()
		if !exists {
			t.Errorf("Expected item to exist")
		}

		if item != i {

			t.Errorf("Expected %d, got %d", i, item)
		}
	}

	size = rb.Size()
	if size != 0 {
		t.Errorf("Expected size %d, got %d", 0, size)
	}
}

func TestRingBufferDoubling(t *testing.T) {
	rb := ringbuffer.NewRingBuffer[int](10)
	for i := range 1000 {
		rb.Enqueue(i)
	}

	size := rb.Size()
	if size != 1000 {
		t.Errorf("Expected size %d, got %d", 1000, size)
	}

	for i := range 1000 {
		item, exists := rb.Dequeue()
		if !exists {
			t.Errorf("Expected item to exist")
		}

		if item != i {

			t.Errorf("Expected %d, got %d", i, item)
		}
	}

	size = rb.Size()
	if size != 0 {
		t.Errorf("Expected size %d, got %d", 0, size)
	}
}

func TestConcurrentEnqueueDequeue(t *testing.T) {
	rb := ringbuffer.NewRingBuffer[int](10)

	enqueue := 0

	dequeue := 0

	iterationCount := 100
	var wg sync.WaitGroup
	wg.Add(2 * iterationCount)

	go func() {
		for range iterationCount {
			rb.Enqueue(enqueue)
			enqueue++
			wg.Done()
		}
	}()

	go func() {
		for dequeue != iterationCount {
			item, exists := rb.Dequeue()
			if !exists {
				continue
			}

			if item != dequeue {
				t.Errorf("Expected %d got %d", dequeue, item)
			}
			dequeue++
			wg.Done()
		}
	}()

	wg.Wait()

	size := rb.Size()
	if size != 0 {
		t.Errorf("Expected size %d, got %d", 0, size)
	}
}
