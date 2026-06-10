package main

import (
	"fmt"
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func main() {
	server := http.NewServeMux()

	server.HandleFunc("/", Handler)

	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatal(err)
	}
}
