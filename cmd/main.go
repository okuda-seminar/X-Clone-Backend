package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"x-clone-backend/api/handlers"
	"x-clone-backend/db"
)

const (
	port = 80
)

func main() {
	db, err := db.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World\n")
	})

	http.HandleFunc("DELETE /api/posts/{postId}", func(w http.ResponseWriter, r *http.Request) {
		postId := r.PathValue("postId")
		fmt.Fprintf(w, "Received Delete request for post id: %s.\n", postId)
		slog.Info(fmt.Sprintf("DELETE /api/posts was called with %s.", postId))
	})

	http.HandleFunc("POST /api/posts", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreatePost(w, r, db)
	})

	http.HandleFunc("POST /api/posts/repost", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Received Post request for reposts.\n")
		slog.Info("POST /api/posts/repost was called.")
	})

	http.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateUser(w, r, db)
	})

	http.HandleFunc("DELETE /api/users/{userID}", func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("userID")
		fmt.Fprintf(w, "Received Delete request for user id: %s.\n", userID)
		slog.Info(fmt.Sprintf("DELETE /api/users was called with %s.", userID))
	})

	http.HandleFunc("GET /api/users/{userID}", func(w http.ResponseWriter, r *http.Request) {
		handlers.FindUserByID(w, r, db)
		slog.Info("GET /api/users/{userID} was called.")
	})

	http.HandleFunc("POST /api/users/{following_user_id}/follow", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Follow a user.")
	})

	http.HandleFunc("DELETE /api/users/{following_user_id}/follow/{followed_user_id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Unfollow a user.")
	})

	http.HandleFunc("POST /api/users/{id}/muting", func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("id")
		fmt.Fprintf(w, "Received POST request for user id: %s.\n", userID)
	})

	http.HandleFunc("DELETE /api/users/{source_user_id}/muting/{target_user_id}", func(w http.ResponseWriter, r *http.Request) {
		sourceUserID := r.PathValue("source_user_id")
		targetUserID := r.PathValue("target_user_id")
		fmt.Fprintf(w, "Received DELETE request for source user id: %s and target user id: %s.\n", sourceUserID, targetUserID)
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
