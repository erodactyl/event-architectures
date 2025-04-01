package main

import (
	sdk "eventbus/sdk/eventbus"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	log.Printf("Received event %s", string(body))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received"))
}

func main() {
	client := sdk.NewEventBusClient("localhost", ":8000")

	port := os.Getenv("PORT")

	callbackURL := fmt.Sprintf("http://localhost:%s/webhook", port)

	unsub := client.Subscribe("Messages", callbackURL)

	http.HandleFunc("/webhook", WebhookHandler)

	fmt.Println("Listening on port", port)
	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		unsub()
		log.Fatal(err)
	}
}
