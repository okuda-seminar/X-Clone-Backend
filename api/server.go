package api

import (
	"database/sql"
	"net/http"
	"x-clone-backend/api/handlers"
	"x-clone-backend/infrastructure"
	"x-clone-backend/usecases"
)

type Server struct {
	db *sql.DB
}

func NewServer(db *sql.DB) Server {
	return Server{db}
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	usersRepository := infrastructure.NewUsersRepository(s.db)
	createUserUsecase := usecases.NewCreateUserUsecase(usersRepository)
	handlers.CreateUser(w, r, createUserUsecase)
}
