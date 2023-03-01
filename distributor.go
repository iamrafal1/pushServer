package main

import (
	"fmt"
	"log"
	"net/http"
)

type MessageChan chan []byte

// A distributor holds open client connections,
// listens for incoming events on its messages channel
// and broadcast event data to all registered connections
type Distributor struct {
	messages       chan string          // Channel with messages
	newClients     chan MessageChan     // New client connections
	closingClients chan MessageChan     // Closed client connections
	clients        map[MessageChan]bool // Client connection map
}

// Listen on various channels. This must run in a go routine
func (d *Distributor) listen() {
	for {
		select {
		case s := <-d.newClients:

			// New client connected, add them to client map
			d.clients[s] = true
			log.Printf("Client added. %d registered clients", len(d.clients))
		case s := <-d.closingClients:

			// Client disconnected, stop sending them messages.
			delete(d.clients, s)
			log.Printf("Removed client. %d registered clients", len(d.clients))
		case messages := <-d.messages:

			// New message detected, broadcase message to all clients
			byte_array := []byte(messages)
			for s := range d.clients {
				s <- byte_array
			}
			log.Printf("Broadcast message to %d clients", len(d.clients))
		}
	}

}

// http handler interface. Each individual connection is handled here
func (d *Distributor) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	flusher, ok := rw.(http.Flusher)

	if !ok {
		http.Error(rw, "Flushing impossible!", http.StatusInternalServerError)
		return
	}

	// Set event streaming headers.
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// Each connection creates its own message channel
	messageChan := make(MessageChan)

	// Notify distributor that new client is created
	d.newClients <- messageChan

	// Remove this client from client map when this handler exits.
	defer func() {
		d.closingClients <- messageChan
	}()

	// Listen to connection close and un-register messageChan
	notify := ctx.Done()

	go func() {
		<-notify
		d.closingClients <- messageChan
	}()

	// block waiting or messages broadcast on this connection's messageChan
	for {
		// Write to the ResponseWriter
		fmt.Fprintf(rw, "data: %s\n\n", <-messageChan)
		// Flush the data immediately instead of buffering it for later.
		flusher.Flush()
	}
}
