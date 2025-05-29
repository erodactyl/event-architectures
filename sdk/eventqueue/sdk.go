package sdk

import (
	"eventarch/pkg/eventqueue"
	"log"
	"net/rpc"
)

type EventQueueClient struct {
	client *rpc.Client
}

func NewEventQueueClient(address, port string) *EventQueueClient {
	client, err := rpc.DialHTTP("tcp", address+port)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	return &EventQueueClient{client}
}

func (c *EventQueueClient) CreateQueue(name string) bool {
	args := &eventqueue.CreateQueueArgs{Name: name}
	var created bool
	err := c.client.Call("EventQueueService.CreateQueue", args, &created)
	if err != nil {
		log.Fatal("EventBusService error: ", err)
	}

	log.Printf("Created queue %s", name)

	return created
}

func (c *EventQueueClient) Put(queueName, eventType string, body []byte) bool {
	args := &eventqueue.PutArgs{QueueName: queueName, EventType: eventType, Body: body}
	var put bool
	err := c.client.Call("EventQueueService.Put", args, &put)
	if err != nil {
		log.Fatal("EventBusService error: ", err)
	}

	log.Printf("Put event %s with type %s in queue %s", string(body), queueName, eventType)
	return put
}

func (c *EventQueueClient) Pull(name string, limit int) []eventqueue.Event {
	args := &eventqueue.PullArgs{QueueName: name, Limit: limit}
	var results []eventqueue.Event
	err := c.client.Call("EventQueueService.Pull", args, &results)
	if err != nil {
		log.Fatal("EventQueueService error: ", err)
	}

	log.Printf("Pulled %d events from %s queue", len(results), name)
	return results
}
