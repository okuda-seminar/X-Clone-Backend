package api

import (
	"database/sql"
	"net/http"
	"x-clone-backend/api/handlers"
)

type Server struct {
	db *sql.DB
}

func NewServer(db *sql.DB) Server {
	return Server{db}
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	handlers.CreateUser(w, r, s.db)
}
