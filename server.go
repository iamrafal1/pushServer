package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/iamrafal1/pushServer/db"
)

func main() {

	data, err := db.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	distributorMap := initialiseDistributors(data)
	fmt.Println(distributorMap)
	http.HandleFunc("/top", WebhookHandler(distributorMap, data))
	http.HandleFunc("/generate", generationHandler(distributorMap, data))
	http.HandleFunc("/delete", deletionHandler(distributorMap, data))
	http.HandleFunc("/", handler)
	http.HandleFunc("/test", handler2)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initialiseDistributors(d *db.Database) map[string]*Distributor {

	// Get urls from database
	urls, err := d.GetAllUrls()
	if err != nil {
		log.Fatal(err)
	}

	// Create url - dist map
	urlMap := make(map[string]*Distributor)

	// Iterate through distributors and initialise them
	for _, url := range urls {
		dist := newDistributor()
		log.Print(url)
		http.Handle("/"+url, dist)
		urlMap[url] = dist
	}

	return urlMap
}

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

func handler2(w http.ResponseWriter, r *http.Request) {

	// Read in the template with SSE JavaScript code.
	t, err := template.ParseFiles("templates/test.html")
	if err != nil {
		log.Fatal("error parsing template.")

	}

	// Render the template, writing to `w`.
	t.Execute(w, nil)

	// Done.
	log.Println("Finished HTTP request at", r.URL.Path)
}

func generationHandler(dists map[string]*Distributor, data *db.Database) func(http.ResponseWriter, *http.Request) {
	if dists == nil {
		log.Fatal("ERROR nil Distributor session!")
	}
	if data == nil {
		log.Fatal("ERROR nil Database session!")
	}

	// Handler for incoming webhooks
	return func(w http.ResponseWriter, r *http.Request) {

		for {

			// Generate values
			h1 := hashGenerator()
			key := h1[0:9]
			h2 := hashGenerator()
			url := h2[0:9]
			h3 := hashGenerator()
			token := h3[0:31]

			// Try to insert into db. If fails, generate again
			_, err := data.InsertAllCols(key, url, token)
			if err != nil {
				continue
			}

			// TODO: change localhost to some url
			// Else, write as json as a response and create distributor
			dist := newDistributor()
			dists[url] = dist
			fmt.Fprintf(w, `{"key": "%s", "url": "127.0.0.1:8080/%s", "token":"%s"}`, key, url, token)
			break
		}

		// Done.
		log.Println("Finished HTTP request at", r.URL.Path)
	}
}

func deletionHandler(dists map[string]*Distributor, data *db.Database) func(http.ResponseWriter, *http.Request) {
	if dists == nil {
		log.Fatal("ERROR nil Distributor session!")
	}
	if data == nil {
		log.Fatal("ERROR nil Database session!")
	}

	// Handler for incoming webhooks
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate key and token
		url := RequestValidator(r, data)
		if url == "" {
			log.Print("Validation failed")
			w.Write([]byte("Validation error"))
			return
		}

		key := r.Header.Get("Push-Key")
		token := r.Header.Get("Push-Token")

		_, err := data.DeleteRow(key, url, token)
		if err != nil {
			log.Print("Failed deletion")
			w.Write([]byte("Error"))
			return
		}
		delete(dists, url)

		w.Write([]byte("Success"))

		// Done.
		log.Println("Finished delete request at", r.URL.Path)
	}
}
