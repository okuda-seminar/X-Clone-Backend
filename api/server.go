package api

import (
	"database/sql"
	"net/http"
	"x-clone-backend/api/handlers"
	"x-clone-backend/internal/app/services"
	"x-clone-backend/internal/app/usecases"
	infrastructure "x-clone-backend/internal/infrastructure/persistence"
)

type Server struct {
	db          *sql.DB
	authService *services.AuthService
}

func NewServer(db *sql.DB, authService *services.AuthService) Server {
	return Server{
		db:          db,
		authService: authService,
	}
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	usersRepository := infrastructure.NewUsersRepository(s.db)
	createUserUsecase := usecases.NewCreateUserUsecase(usersRepository)
	handlers.CreateUser(w, r, createUserUsecase, s.authService)
}
