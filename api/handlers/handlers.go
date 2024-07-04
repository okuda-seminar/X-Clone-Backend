package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"x-clone-backend/entities"

	"github.com/google/uuid"
)

// CreateUser creates a new user with the specified useranme and display name,
// then, inserts it into users table.
func CreateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var body createUserRequestBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintln("Request body was invalid."), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO users (id, username, display_name, bio, is_private) VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at`

	var createdAt, updatedAt time.Time
	id := uuid.New()

	err = db.QueryRow(query, id, body.Username, body.DisplayName, "", false).Scan(&createdAt, &updatedAt)
	if err != nil {
		http.Error(w, fmt.Sprintln("Could not create a user."), http.StatusInternalServerError)
		return
	}

	user := entities.User{
		ID:          id,
		Username:    body.Username,
		DisplayName: body.DisplayName,
		Bio:         "",
		IsPrivate:   false,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)
	err = encoder.Encode(&user)
	if err != nil {
		http.Error(w, fmt.Sprintln("Could not encode response."), http.StatusInternalServerError)
		return
	}
}
