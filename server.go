package main

import (
	"bytes"
	"encoding/json"
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
	http.HandleFunc("/top/", webhookHandler(distributor))
	go func() {
		for i := 0; ; i++ {
			time.Sleep(5e9)
			//Encode the data
			postBody, _ := json.Marshal(map[string]string{
				"name":  "Hello",
				"email": "Friends@otz.com",
			})
			responseBody := bytes.NewBuffer(postBody)
			//Leverage Go's HTTP Post function to make request
			_, err := http.Post("http://127.0.0.1:8080/top/", "application/json", responseBody)
			if err != nil {
				log.Printf("Something really bad happened")
				log.Printf(err.Error())
			}
			// Print log message and sleep for 5 seconds.
			log.Printf("Sent message %d ", i)
		}
	}()
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
