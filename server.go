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
	http.HandleFunc("/top/", WebhookHandler(distributorMap, data))
	http.HandleFunc("/", handler)
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
