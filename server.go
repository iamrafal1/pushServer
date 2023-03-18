package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
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
	go dbTest()
	http.Handle("/events/", distributor)
	http.HandleFunc("/top/", webhookHandler(distributor))
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func dbTest() {
	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	var version string
	err = db.QueryRow("SELECT SQLITE_VERSION()").Scan(&version)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(version)
}
