package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"x-clone-backend/api"
	"x-clone-backend/api/handlers"
	"x-clone-backend/api/middlewares"
	"x-clone-backend/db"
	openapi "x-clone-backend/gen"
	"x-clone-backend/internal/app/usecases"
	infrastructure "x-clone-backend/internal/infrastructure/persistence"
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

	var userChannels = make(map[string]chan []byte)
	var mu sync.Mutex

	server := api.NewServer(db)
	mux := http.NewServeMux()
	postsRepository := infrastructure.NewPostsRepository(db)
	getSpecificUserPostsUsecase := usecases.NewGetSpecificUserPostsUsecase(postsRepository)
	getUserAndFolloweePostsUsecase := usecases.NewGetUserAndFolloweePostsUsecase(postsRepository)

	usersRepository := infrastructure.NewUsersRepository(db)
	deleteUserUsecase := usecases.NewDeleteUserUsecase(usersRepository)
	getspecificUserUsecase := usecases.NewGetSpecificUserUsecase(usersRepository)
	likePostUsecase := usecases.NewLikePostUsecase(usersRepository)
	unlikePostUsecase := usecases.NewUnlikePostUsecase(usersRepository)
	followUserUsecase := usecases.NewFollowUserUsecase(usersRepository)
	unfollowUserUsecase := usecases.NewUnfollowUserUsecase(usersRepository)
	muteUserUsecase := usecases.NewMuteUserUsecase(usersRepository)
	unmuteUserUsecase := usecases.NewUnmuteUserUsecase(usersRepository)
	blockUserUsecase := usecases.NewBlockUserUsecase(usersRepository)
	unblockUserUsecase := usecases.NewUnblockUserUsecase(usersRepository)

	mux.HandleFunc("POST /api/posts", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreatePost(w, r, db, &mu, &userChannels)
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
		handlers.CreateFollowship(w, r, followUserUsecase)
	})

	mux.HandleFunc("GET /api/users/{id}/posts", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUserPostsTimeline(w, r, getSpecificUserPostsUsecase)
	})

	mux.HandleFunc("GET /api/users/{id}/timelines/reverse_chronological", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetReverseChronologicalHomeTimeline(w, r, getUserAndFolloweePostsUsecase, &mu, &userChannels)
	})

	mux.HandleFunc("DELETE /api/users/{source_user_id}/following/{target_user_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteFollowship(w, r, unfollowUserUsecase)
	})

	mux.HandleFunc("POST /api/users/{id}/muting", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateMuting(w, r, muteUserUsecase)
	})

	mux.HandleFunc("DELETE /api/users/{source_user_id}/muting/{target_user_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteMuting(w, r, unmuteUserUsecase)
	})

	mux.HandleFunc("POST /api/users/{id}/blocking", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateBlocking(w, r, blockUserUsecase)
	})

	mux.HandleFunc("DELETE /api/users/{source_user_id}/blocking/{target_user_id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteBlocking(w, r, unblockUserUsecase)
	})

	mux.HandleFunc("/api/notifications", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Notifications\n")
	})

	handler := middlewares.CORS(openapi.HandlerFromMux(&server, mux))
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
