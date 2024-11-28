package api

import (
	"database/sql"
	"log"
	"net/http"
	"x-clone-backend/api/handlers"
	"x-clone-backend/internal/app/usecases"
	infrastructure "x-clone-backend/internal/infrastructure/persistence"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	db *sql.DB
}

func NewServer(db *sql.DB) Server {
	return Server{db}
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: s.db,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to wrap *sql.DB with GORM: %v", err)
	}

	usersRepository := infrastructure.NewUsersRepository(gormDB)
	createUserUsecase := usecases.NewCreateUserUsecase(usersRepository)
	handlers.CreateUser(w, r, createUserUsecase)
}
