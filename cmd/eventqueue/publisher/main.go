package main

import (
	"encoding/json"
	sdk "eventarch/sdk/eventqueue"
	"time"
)

func main() {
	client := sdk.NewEventQueueClient("localhost", ":8000")

	client.CreateQueue("default")

	for range 10 {
		data := `{"hello": "world"}`
		client.Put("default", "default_type", json.RawMessage(data))
		time.Sleep(1 * time.Second)
	}
}
