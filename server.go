package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/iamrafal1/pushServer/db"
)

func handler(w http.ResponseWriter, r *http.Request) {

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
	data, err := db.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	fmt.Println(data.GetAllUrls())
	http.Handle("/events/", distributor)
	http.HandleFunc("/top/", webhookHandler(distributor))
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
