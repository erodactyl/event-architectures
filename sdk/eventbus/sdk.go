package sdk

import (
	"eventbus/pkg/eventbus"
	"log"
	"net/rpc"
)

type EventBusClient struct {
	client *rpc.Client
}

func NewEventBusClient(address, port string) *EventBusClient {
	client, err := rpc.DialHTTP("tcp", address+port)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	return &EventBusClient{client}
}

func (c *EventBusClient) Subscribe(busName, callbackURL string) func() bool {
	args := &eventbus.SubscribeArgs{BusName: busName, CallbackURL: callbackURL}
	var id string
	err := c.client.Call("EventBusService.Subscribe", args, &id)
	if err != nil {
		log.Fatal("EventBusService error: ", err)
	}

	log.Printf("Subscribed to event bus %s with callback url %s", busName, callbackURL)

	return func() bool {
		return c.unsubscribe(busName, id)
	}
}

func (c *EventBusClient) Publish(busName string, eventType string, body []byte) {
	args := &eventbus.PublishArgs{BusName: busName, EventType: eventType, Body: body}
	var succeess bool
	err := c.client.Call("EventBusService.Publish", args, &succeess)
	if err != nil {
		log.Fatal("EventBusService error: ", err)
	}

	log.Printf("Published event %s to bus %s", string(body), busName)
}

func (c *EventBusClient) unsubscribe(busName, id string) bool {
	args := &eventbus.UnsubscribeArgs{BusName: busName, ID: id}
	var unsubscribed bool
	err := c.client.Call("EventBusService.Unsubscribe", args, &unsubscribed)
	if err != nil {
		log.Fatal("EventBusService error: ", err)
	}

	log.Printf("Unsubscribed from bus %s", busName)

	return unsubscribed
}
