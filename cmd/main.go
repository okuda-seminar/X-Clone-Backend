package main

import (
	"fmt"
	"log"
	"net/http"
	"x-clone-backend/db"
)

const (
	port = 80
)

func main() {
	_, err := db.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World")
	})

	log.Println("Starting server...")

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalln(err)
	}
}
