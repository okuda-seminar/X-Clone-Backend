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

	query := `INSERT INTO users (username, display_name, bio, is_private) VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	var (
		id                   uuid.UUID
		createdAt, updatedAt time.Time
	)

	err = db.QueryRow(query, body.Username, body.DisplayName, "", false).Scan(&id, &createdAt, &updatedAt)
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

// FindUserByID finds a user with the specified ID.
func FindUserByID(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID := r.PathValue("userID")

	query := `SELECT * FROM users WHERE id = $1`
	row := db.QueryRow(query, userID)

	var user entities.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.DisplayName,
		&user.Bio,
		&user.IsPrivate,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not find a user (ID: %s)\n", userID), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	err = encoder.Encode(&user)
	if err != nil {
		http.Error(w, fmt.Sprintln("Could not encode response."), http.StatusInternalServerError)
		return
	}
}

// CreatePost creates a new post with the specified user_id and text,
// then, inserts it into posts table.
func CreatePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var body createPostRequestBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintln("Request body was invalid."), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO posts (user_id, text) VALUES ($1, $2)
		RETURNING id, created_at`

	var (
		id        uuid.UUID
		createdAt time.Time
	)

	err = db.QueryRow(query, body.UserID, body.Text).Scan(&id, &createdAt)
	if err != nil {
		http.Error(w, fmt.Sprintln("Could not create a post."), http.StatusInternalServerError)
		return
	}

	post := entities.Post{
		ID:        id,
		UserID:    body.UserID,
		Text:      body.Text,
		CreatedAt: createdAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)
	err = encoder.Encode(&post)
	if err != nil {
		http.Error(w, fmt.Sprintln("Could not encode response."), http.StatusInternalServerError)
		return
	}
}
