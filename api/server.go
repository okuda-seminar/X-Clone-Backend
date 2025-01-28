package api

import (
	"database/sql"
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
	handlers.CreateRepostHandler
	handlers.DeleteRepostHandler
	handlers.GetUserPostsTimelineHandler
	handlers.GetReverseChronologicalHomeTimelineHandler
}

func NewServer(db *sql.DB, mu *sync.Mutex, usersChan *map[string]chan entities.TimelineEvent) Server {
	return Server{
		CreateUserHandler:                          handlers.NewCreateUserHandler(db),
		FindUserByIDHandler:                        handlers.NewFindUserByIDHandler(db),
		CreatePostHandler:                          handlers.NewCreatePostHandler(db, mu, usersChan),
		CreateRepostHandler:                        handlers.NewCreateRepostHandler(db, mu, usersChan),
		DeleteRepostHandler:                        handlers.NewDeleteRepostHandler(db),
		GetUserPostsTimelineHandler:                handlers.NewGetUserPostsTimelineHandler(db),
		GetReverseChronologicalHomeTimelineHandler: handlers.NewGetReverseChronologicalHomeTimelineHandler(db, mu, usersChan),
	}
}
