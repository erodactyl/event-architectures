package main

import (
	"encoding/json"
	sdk "eventbus/sdk/eventbus"
	"time"
)

func main() {
	client := sdk.NewEventBusClient("localhost", ":8000")

	for range 10 {
		data := `{"hello": "world"}`
		client.Publish("Messages", "main", json.RawMessage(data))
		time.Sleep(1 * time.Second)
	}
}
