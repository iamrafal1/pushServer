package handlers

import (
	"fmt"
	"log"
	"net/http"
)

// Observer-like data structure
type Distributor struct {
	messages chan string          // Channel for messages from the outside
	clients  map[chan string]bool // Client connection map
}

// Creates new distributor instances and makes them listen in a go routine
func NewDistributor() *Distributor {
	dist := &Distributor{
		messages: make(chan string),
		clients:  make(map[chan string]bool),
	}
	go dist.listen()
	return dist
}

// Listen for new message update. This must run in a go routine
func (d *Distributor) listen() {
	for {
		select {
		// New message detected, broadcast message to all clients
		case messages := <-d.messages:
			for s := range d.clients {
				s <- messages
			}
			log.Printf("Broadcast message to %d clients", len(d.clients))
		}
	}

}

// http handler interface. Each individual connection is handled here
func (d *Distributor) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// Get request context and create a flusher
	ctx := req.Context()
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Flushing impossible!", http.StatusInternalServerError)
		return
	}

	// Set event streaming headers.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Each SSE connection creates its own connection
	conn := make(chan string)

	// Add new client to client map
	d.clients[conn] = true
	log.Printf("Client added. %d registered clients", len(d.clients))

	// Listen to connection close and un-register connection
	notify := ctx.Done()
	go func() {
		<-notify
		delete(d.clients, conn)
		log.Printf("Removed client. %d registered clients", len(d.clients))
	}()

	// Loop infinitely, flush messages as they arrive
	for {
		// Write to the ResponseWriter and flush
		fmt.Fprintf(w, "data: %s\n\n", <-conn)
		flusher.Flush()
	}
}
