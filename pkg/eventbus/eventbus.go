package eventbus

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

type subscriber struct {
	ID          string
	BusName     string
	CallbackURL string
}

type Event struct {
	// Name of the event bus to send the event to
	BusName string

	EventType string

	// JSON body of the event
	Body []byte
}

type EventBusService struct {
	subscribers map[string][]*subscriber
	mu          sync.RWMutex
}

func NewEventBusService() *EventBusService {
	return &EventBusService{
		subscribers: map[string][]*subscriber{},
		mu:          sync.RWMutex{},
	}
}

func (e *EventBusService) handleEvent(event *Event) {
	e.mu.RLock()
	subs, exists := e.subscribers[event.BusName]
	if !exists {
		e.mu.RUnlock()
		return
	}

	urls := make([]string, len(subs))
	for i, sub := range subs {
		urls[i] = sub.CallbackURL
	}
	e.mu.RUnlock()

	var wg sync.WaitGroup

	success := 0
	fail := 0
	var mu sync.Mutex // Protect reply counters

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			body := bytes.NewReader(event.Body)
			req, err := http.NewRequest(http.MethodPost, url, body)
			if err != nil {
				mu.Lock()
				fail++
				mu.Unlock()
				return
			}

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				mu.Lock()
				fail++
				mu.Unlock()
				return
			}
			defer res.Body.Close()

			if res.StatusCode == http.StatusOK {
				mu.Lock()
				success++
				mu.Unlock()
			} else {
				mu.Lock()
				fail++
				mu.Unlock()
			}
		}(url)
	}

	wg.Wait()

	log.Printf("Processed event for queue %s", event.BusName)
}

type PublishArgs = Event

func (e *EventBusService) Publish(event *PublishArgs, reply *bool) error {
	go e.handleEvent(event)
	*reply = true
	return nil
}

type SubscribeArgs struct {
	BusName     string
	CallbackURL string
}

func (e *EventBusService) Subscribe(args *SubscribeArgs, reply *string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	sub := &subscriber{uuid.New().String(), args.BusName, args.CallbackURL}

	if _, exists := e.subscribers[args.BusName]; exists {
		for _, sub := range e.subscribers[args.BusName] {
			if sub.CallbackURL == args.CallbackURL {
				return errors.New("Subscriber already exists")
			}
		}
		e.subscribers[args.BusName] = append(e.subscribers[args.BusName], sub)
	} else {
		e.subscribers[args.BusName] = []*subscriber{sub}
	}

	*reply = sub.ID

	log.Printf("Registered subscriber for bus %s with callback URL %s", args.BusName, args.CallbackURL)

	return nil
}

// O(1) remove an element from a slice by replacing it with the last element
func remove[K any](s []K, index int) []K {
	s[index] = s[len(s)-1]
	return s[:len(s)-1]
}

type UnsubscribeArgs struct {
	ID      string
	BusName string
}

func (e *EventBusService) Unsubscribe(args *UnsubscribeArgs, reply *bool) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.subscribers[args.BusName]; !exists {
		return errors.New(fmt.Sprintf("Bus with name %s does not exist", args.BusName))
	}

	for i, sub := range e.subscribers[args.BusName] {
		if sub.ID == args.ID {
			e.subscribers[args.BusName] = remove(e.subscribers[args.BusName], i)
			*reply = true
			return nil
		}
	}

	return errors.New(fmt.Sprintf("Subscriber with ID %s for bus %s does not exist", args.ID, args.BusName))
}
