package main

import (
	"eventarch/pkg/eventbus"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	eventBusService := eventbus.NewEventBusService()
	rpc.Register(eventBusService)
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
