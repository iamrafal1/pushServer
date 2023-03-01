package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Read in the template with SSE JavaScript code.
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal("error parsing template.")

	}

	// Render the template, writing to `w`.
	t.Execute(w, "friends")

	// Done.
	log.Println("Finished HTTP request at", r.URL.Path)
}

func main() {
	distributor := &Distributor{
		messages:       make(chan string),
		newClients:     make(chan MessageChan),
		closingClients: make(chan MessageChan),
		clients:        make(map[MessageChan]bool),
	}
	go distributor.listen()
	http.Handle("/events/", distributor)

	go func() {
		for i := 0; ; i++ {

			// Create a message for clients
			distributor.messages <- fmt.Sprintf("%d - the time is %v", i, time.Now())

			// Print log message and sleep for 5 seconds.
			log.Printf("Sent message %d ", i)
			time.Sleep(5e9)

		}
	}()
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
