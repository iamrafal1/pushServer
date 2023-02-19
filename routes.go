package main

import (
	"net/http"
)

func getRoutes() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/handler", handler2)
}
