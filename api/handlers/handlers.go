package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
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

// DeleteUser deletes a user with the specified user ID.
// If a target user does not exist, it returns 404.
func DeleteUserByID(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID := r.PathValue("userID")

	slog.Info(fmt.Sprintf("DELETE /api/users was called with %s.", userID))

	query := `DELETE FROM users WHERE id = $1`
	res, err := db.Exec(query, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not delete a user (ID: %s)\n", userID), http.StatusInternalServerError)
		return
	}
	count, err := res.RowsAffected()
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not delete a user (ID: %s)\n", userID), http.StatusInternalServerError)
		return
	}
	if count != 1 {
		http.Error(w, fmt.Sprintf("No row found to delete (ID: %s)\n", userID), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

func CreateFollowship(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var body createFollowshipRequestBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintln("Request body was invalid."), http.StatusBadRequest)
		return
	}

	sourceUserID := r.PathValue("id")

	query := `INSERT INTO followships (source_user_id, target_user_id) VALUES ($1, $2)`

	_, err = db.Exec(query, sourceUserID, body.TargetUserID)
	if err != nil {
		http.Error(w, fmt.Sprintln("Could not create followship."), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteFollowship(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	sourceUserID := r.PathValue("source_user_id")
	targetUserID := r.PathValue("target_user_id")

	query := `DELETE FROM followships WHERE source_user_id = $1 AND target_user_id = $2`
	res, err := db.Exec(query, sourceUserID, targetUserID)
	if err != nil {
		http.Error(w, "Could not delete followship.", http.StatusInternalServerError)
		return
	}
	count, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Could not delete followship.", http.StatusInternalServerError)
		return
	}
	if count != 1 {
		http.Error(w, "No row found to delete.", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CreateRepost creates a new repost with the specified post_id and user_id,
// then, inserts it into reposts table.
func CreateRepost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var body createRepostRequestBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)

	if err != nil {
		http.Error(w, fmt.Sprintln("Request body was invalid."), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO reposts (post_id, user_id) VALUES ($1, $2)`

	_, err = db.Exec(query, body.PostID, body.UserID)

	if err != nil {
		http.Error(w, fmt.Sprintln("Could not create a repost."), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
