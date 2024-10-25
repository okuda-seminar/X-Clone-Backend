package main

import (
	"fmt"
	"log"
	"net/http"
	"x-clone-backend/api"
	"x-clone-backend/api/handlers"
	"x-clone-backend/api/middlewares"
	"x-clone-backend/db"
	openapi "x-clone-backend/gen"
	"x-clone-backend/infrastructure"
	"x-clone-backend/usecases"
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

	sever := api.NewServer(db)
	mux := http.NewServeMux()
	postsRepository := infrastructure.NewPostsRepository(db)
	getSpecificUserPostsUsecase := usecases.NewGetSpecificUserPostsUsecase(postsRepository)
	getUserAndFolloweePostsUsecase := usecases.NewGetUserAndFolloweePostsUsecase(postsRepository)

	usersRepository := infrastructure.NewUsersRepository(db)
	deleteUserUsecase := usecases.NewDeleteUserUsecase(usersRepository)
	getspecificUserUsecase := usecases.NewGetSpecificUserUsecase(usersRepository)
	likePostUsecase := usecases.NewLikePostUsecase(usersRepository)
	unlikePostUsecase := usecases.NewUnlikePostUsecase(usersRepository)

	mux.HandleFunc("POST /api/posts", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreatePost(w, r, db)
	})

	mux.HandleFunc("DELETE /api/posts/{postID}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeletePost(w, r, db)
	})

	mux.HandleFunc("POST /api/posts/reposts", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateRepost(w, r, db)
	})

	mux.HandleFunc("DELETE /api/posts/reposts/{user_id}/{post_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteRepost(w, r, db)
	})

	mux.HandleFunc("DELETE /api/users/{userID}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteUserByID(w, r, deleteUserUsecase)
	})

	mux.HandleFunc("GET /api/users/{userID}", func(w http.ResponseWriter, r *http.Request) {
		handlers.FindUserByID(w, r, getspecificUserUsecase)
	})

	mux.HandleFunc("POST /api/users/{id}/likes", func(w http.ResponseWriter, r *http.Request) {
		handlers.LikePost(w, r, likePostUsecase)
	})

	mux.HandleFunc("DELETE /api/users/{id}/likes/{post_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.UnlikePost(w, r, unlikePostUsecase)
	})

	mux.HandleFunc("POST /api/users/{id}/following", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateFollowship(w, r, db)
	})

	mux.HandleFunc("GET /api/users/{id}/posts", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUserPostsTimeline(w, r, getSpecificUserPostsUsecase)
	})

	mux.HandleFunc("GET /api/users/{id}/timelines/reverse_chronological", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetReverseChronologicalHomeTimeline(w, r, getUserAndFolloweePostsUsecase)
	})

	mux.HandleFunc("DELETE /api/users/{source_user_id}/following/{target_user_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteFollowship(w, r, db)
	})

	mux.HandleFunc("POST /api/users/{id}/muting", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateMuting(w, r, db)
	})

	mux.HandleFunc("DELETE /api/users/{source_user_id}/muting/{target_user_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteMuting(w, r, db)
	})

	mux.HandleFunc("POST /api/users/{id}/blocking", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateBlocking(w, r, db)
	})

	mux.HandleFunc("DELETE /api/users/{source_user_id}/blocking/{target_user_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteBlocking(w, r, db)
	})

	mux.HandleFunc("/api/notifications", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Notifications\n")
	})

	handler := middlewares.CORS(openapi.HandlerFromMux(&sever, mux))
	s := http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%d", port),
	}

	log.Println("Starting server...")

	err = s.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
