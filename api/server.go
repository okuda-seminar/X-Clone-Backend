package api

import (
	"database/sql"
	"net/http"
	"x-clone-backend/api/handlers"
)

// [Server] satisfies [ServerInterface] defined in gen/server.gen.go.
type Server struct {
	handlers.CreateUserHandler
}

func NewServer(db *sql.DB) Server {
	return Server{
		CreateUserHandler: handlers.NewCreateUserHandler(db),
	}
}

// Define temporary handlers so that [Server] satisfies [ServerInterface].
func (s *Server) CreatePost(w http.ResponseWriter, r *http.Request)                                 {}
func (s *Server) CreateRepost(w http.ResponseWriter, r *http.Request)                               {}
func (s *Server) DeleteRepost(w http.ResponseWriter, r *http.Request, userId string, postId string) {}
func (s *Server) GetUserPostsTimeline(w http.ResponseWriter, r *http.Request, id string)            {}
func (s *Server) GetReverseChronologicalHomeTimeline(w http.ResponseWriter, r *http.Request, id string) {
}
func (s *Server) FindUserByID(w http.ResponseWriter, r *http.Request, userID string) {}
