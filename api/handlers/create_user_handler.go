package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"x-clone-backend/api/transfers"
	openapi "x-clone-backend/gen"
	"x-clone-backend/internal/app/usecases"
	infrastructure "x-clone-backend/internal/infrastructure/persistence"
)

type CreateUserHandler struct {
	createUserUsecase usecases.CreateUserUsecase
}

func NewCreateUserHandler(db *sql.DB) CreateUserHandler {
	usersRepository := infrastructure.NewUsersRepository(db)
	createUserUsecase := usecases.NewCreateUserUsecase(usersRepository)
	return CreateUserHandler{
		createUserUsecase,
	}
}

// CreateUser creates a new user with the specified user infomtaion.
func (h *CreateUserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var body openapi.CreateUserRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintln("Request body was invalid."), http.StatusBadRequest)
		return
	}

	user, err := h.createUserUsecase.CreateUser(body.Username, body.DisplayName, body.Password)
	if err != nil {
		var code int

		if isUniqueViolationError(err) {
			code = http.StatusConflict
		} else {
			code = http.StatusInternalServerError
		}
		http.Error(w, fmt.Sprintln("Could not create a user."), code)
		return
	}
	res := transfers.ToCreateUserResponse(&user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)
	err = encoder.Encode(res)
	if err != nil {
		http.Error(w, fmt.Sprintln("Could not encode response."), http.StatusInternalServerError)
		return
	}
}
