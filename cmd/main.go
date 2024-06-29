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
		fmt.Fprintf(w, "Hello, World\n")
	})

	http.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			fmt.Fprintf(w, "Received Post request for posts.\n")
		case http.MethodDelete:
			fmt.Fprintf(w, "Received Delete request for posts.\n")
		default:
			http.Error(w, fmt.Sprintln("/api/posts supports only Post and Delete now."), http.StatusHTTPVersionNotSupported)
		}
	})

	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Users\n")
	})

	http.HandleFunc("/api/notifications", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Notifications\n")
	})

	log.Println("Starting server...")

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalln(err)
	}
}
