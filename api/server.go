package api

import (
	"database/sql"
	"net/http"
	"sync"
	"x-clone-backend/api/handlers"
	openapi "x-clone-backend/gen"
	"x-clone-backend/internal/domain/entities"
)

var _ openapi.ServerInterface = (*Server)(nil)

// [Server] satisfies [ServerInterface] defined in gen/server.gen.go.
type Server struct {
	handlers.CreateUserHandler
	handlers.FindUserByIDHandler
	handlers.CreatePostHandler
}

func NewServer(db *sql.DB, mu *sync.Mutex, usersChan *map[string]chan entities.TimelineEvent) Server {
	return Server{
		CreateUserHandler:   handlers.NewCreateUserHandler(db),
		FindUserByIDHandler: handlers.NewFindUserByIDHandler(db),
		CreatePostHandler:   handlers.NewCreatePostHandler(db, mu, usersChan),
	}
}

// Define temporary handlers so that [Server] satisfies [ServerInterface].
func (s *Server) CreateRepost(w http.ResponseWriter, r *http.Request)                               {}
func (s *Server) DeleteRepost(w http.ResponseWriter, r *http.Request, userId string, postId string) {}
func (s *Server) GetUserPostsTimeline(w http.ResponseWriter, r *http.Request, id string)            {}
func (s *Server) GetReverseChronologicalHomeTimeline(w http.ResponseWriter, r *http.Request, id string) {
}
