package main

import (
	sdk "eventarch/sdk/eventqueue"
	"log"
	"time"
)

func main() {
	client := sdk.NewEventQueueClient("localhost", ":8000")

	for {
		results := client.Pull("default", 10)
		for i := 0; i < len(results); i++ {
			log.Printf("Read event of type %s with body %s: ", results[i].EventType, string(results[i].Body))
		}
		time.Sleep(time.Second)
	}
}
