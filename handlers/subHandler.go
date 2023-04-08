package handlers

import (
	"log"
	"net/http"
	"strings"
)

// Wrapper of handler that manages subscriber connections
func SubHandler(dists map[string]*Distributor) func(http.ResponseWriter, *http.Request) {
	if dists == nil {
		log.Fatal("ERROR nil distributor session!")
	}

	// Handler for managing incoming subscriber connections
	return func(w http.ResponseWriter, r *http.Request) {

		// Determine which distributor the subscriber was trying to connect to.
		// This information is contained in the url, after the second /
		path := r.URL.Path
		split := strings.Split(path, "/")
		url := split[2]

		// Retrieve said distributor
		distributor := dists[url]
		if distributor == nil {
			w.Write([]byte("Invalid URL!"))
			log.Print("Invalid URL!")
			return
		}

		// Let the distributor handle the connection from now on.
		// When the distributor connection terminates, simply terminate this connection too
		distributor.ServeHTTP(w, r)
	}
}
