package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "aloha", r.URL.Path[1:])
}

func handler2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Website part 2")
}

func main() {
	getRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
