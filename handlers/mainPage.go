package handlers

import (
	"html/template"
	"log"
	"net/http"
)

// Handles main page of server
func MainHandler(w http.ResponseWriter, r *http.Request) {

	// Read in the template with main webpage
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal("error parsing template.")

	}

	// Render the template
	t.Execute(w, nil)

	// Done.
	log.Println("Finished HTTP request at", r.URL.Path)
}
