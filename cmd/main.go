package main

import (
	"fmt"
	"log"
	"log/slog"
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

	http.HandleFunc("DELETE /api/posts/{postId}", func(w http.ResponseWriter, r *http.Request) {
		postId := r.PathValue("postId")
		fmt.Fprintf(w, "Received Delete request for post id: %s.\n", postId)
		slog.Info(fmt.Sprintf("DELETE /api/posts was called with %s.", postId))
	})

	http.HandleFunc("POST /api/posts", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Received Post request for posts.\n")
		slog.Info("POST /api/posts was called.")
	})

	http.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Create a user.")
	})

	http.HandleFunc("GET /api/users/{userId}", func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("userId")
		fmt.Fprintf(w, "Find a user with the specified ID (%s).\n", userId)
	})

	http.HandleFunc("POST /api/users/{following_user_id}/follow", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Follow a user.")
	})

	http.HandleFunc("DELETE /api/users/{following_user_id}/follow/{followed_user_id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Unfollow a user.")
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
