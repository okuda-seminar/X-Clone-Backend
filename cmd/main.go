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

	http.HandleFunc("POST /api/posts/reposts", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateRepost(w, r, db)
	})

	http.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateUser(w, r, db)
	})

	http.HandleFunc("DELETE /api/users/{userID}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteUserByID(w, r, db)
	})

	http.HandleFunc("GET /api/users/{userID}", func(w http.ResponseWriter, r *http.Request) {
		handlers.FindUserByID(w, r, db)
		slog.Info("GET /api/users/{userID} was called.")
	})

	http.HandleFunc("POST /api/users/{userID}/likes", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Received POST request for likes.\n")
		slog.Info("POST /api/users/{userID}/likes was called.")
	})

	http.HandleFunc("DELETE /api/users/{userID}/likes/{postID}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Received DELETE request for likes\n")
		slog.Info("DELETE /api/users/{userID}/likes/{postID} was called.")

	})

	http.HandleFunc("POST /api/users/{following_user_id}/following", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateFollowship(w, r, db)
	})

	// TODO: https://github.com/okuda-seminar/X-Clone-Backend/issues/54
	// - Rename the endpoint for unfollowing a user in the same way as X does.
	http.HandleFunc("DELETE /api/users/{following_user_id}/following/{followed_user_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteFollowship(w, r, db)
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
