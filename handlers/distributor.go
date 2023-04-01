package handlers

import (
	"fmt"
	"log"
	"net/http"
)

// NOTE (chan string) is essentially like a connection - it's a way to communicate with the go routine that holds the connection. Hence each "connection" below refers to a chan string, but this only applies in relation to clients. Note that messages is not a "connection" despite being the same data type, because logically it's different.
type connection chan string

// Observer-like data structure
type Distributor struct {
	messages       chan string         // Channel for messages from the outside
	newClients     chan connection     // Channel for new client connections
	closingClients chan connection     // Channel for closed client connections
	clients        map[connection]bool // Client connection map
}

func NewDistributor() *Distributor {
	dist := &Distributor{
		messages:       make(chan string),
		newClients:     make(chan connection),
		closingClients: make(chan connection),
		clients:        make(map[connection]bool),
	}
	go dist.listen()
	return dist
}

// Listen on various channels. This must run in a go routine
func (d *Distributor) listen() {
	for {
		select {
		// New message detected, broadcast message to all clients
		case messages := <-d.messages:
			for s := range d.clients {
				s <- messages
			}
			log.Printf("Broadcast message to %d clients", len(d.clients))
		// New client connected, add them to client map
		case conn := <-d.newClients:
			d.clients[conn] = true
			log.Printf("Client added. %d registered clients", len(d.clients))
		// Client disconnected, remove them from the client map
		case conn := <-d.closingClients:
			delete(d.clients, conn)
			log.Printf("Removed client. %d registered clients", len(d.clients))
		}
	}

}

// http handler interface. Each individual connection is handled here
func (d *Distributor) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	// Get request context and create a flusher
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

	// Each SSE connection creates its own communication connection
	communicationConn := make(connection)

	// Notify distributor that new client is created
	d.newClients <- communicationConn

	// Listen to connection close and un-register connection
	notify := ctx.Done()
	go func() {
		<-notify
		d.closingClients <- communicationConn
	}()

	// Loop infinitely, flush messages as they arrive
	for {
		// Write to the ResponseWriter and flush
		fmt.Fprintf(rw, "data: %s\n\n", <-communicationConn)
		flusher.Flush()
	}
}
