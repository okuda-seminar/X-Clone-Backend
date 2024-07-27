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

	http.HandleFunc("POST /api/users/{id}/likes", func(w http.ResponseWriter, r *http.Request) {
		handlers.LikePost(w, r, db)
		slog.Info("POST /api/users/{id}/likes was called.")
	})

	http.HandleFunc("DELETE /api/users/{id}/likes/{post_id}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Received DELETE request for likes\n")
		slog.Info("DELETE /api/users/{id}/likes/{post_id} was called.")

	})

	http.HandleFunc("POST /api/users/{id}/following", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateFollowship(w, r, db)
	})

	http.HandleFunc("DELETE /api/users/{source_user_id}/following/{target_user_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteFollowship(w, r, db)
	})

	http.HandleFunc("POST /api/users/{id}/muting", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateMuting(w, r, db)
	})

	http.HandleFunc("DELETE /api/users/{source_user_id}/muting/{target_user_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteMuting(w, r, db)
	})

	http.HandleFunc("POST /api/users/{id}/blocking", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateBlocking(w, r, db)
	})

	http.HandleFunc("DELETE /api/users/{source_user_id}/blocking/{target_user_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteBlocking(w, r, db)
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
