package main

import (
	"eventarch/pkg/eventqueue"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	eventQueueService := eventqueue.NewEventQueueService()
	rpc.Register(eventQueueService)
	rpc.HandleHTTP()

	port := ":8000"

	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("Listen error:", err)
	}

	log.Print("Starting server on port ", port)
	err = http.Serve(l, nil)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
