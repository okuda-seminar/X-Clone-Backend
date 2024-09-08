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
		var code int

		if isUniqueViolationError(err) {
			code = http.StatusConflict
		} else {
			code = http.StatusInternalServerError
		}
		http.Error(w, fmt.Sprintln("Could not create a user."), code)
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

// DeletePost deletes a post with the specified post ID.
// If the post doesn't exist, it returns 404 error.
func DeletePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	postID := r.PathValue("postID")
	slog.Info(fmt.Sprintf("DELETE /api/posts was called with %s.", postID))

	query := `DELETE FROM posts WHERE id = $1`
	res, err := db.Exec(query, postID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not delete a post (ID: %s)\n", postID), http.StatusInternalServerError)
		return
	}
	count, err := res.RowsAffected()
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not delete a post (ID: %s)\n", postID), http.StatusInternalServerError)
		return
	}
	if count != 1 {
		http.Error(w, fmt.Sprintf("No row found to delete (ID: %s)\n", postID), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// LikePost creates a like with the specified user_id and post_id,
// then, inserts it into likes table.
func LikePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var body likePostRequestBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintln("Request body was invalid."), http.StatusBadRequest)
		return
	}

	userID := r.PathValue("id")

	slog.Info(fmt.Sprintf("POST /api/users/{id}/likes was called with %s.", userID))

	query := "INSERT INTO likes (user_id, post_id) VALUES ($1, $2)"

	_, err = db.Exec(query, userID, body.PostID)
	if err != nil {
		http.Error(w, fmt.Sprintln("Could not create a like."), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UnlikePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID := r.PathValue("id")
	postID := r.PathValue("post_id")

	slog.Info(fmt.Sprintf("DELETE /api/users/{id}/likes/{post_id} was called with %s and %s.", userID, postID))

	query := "DELETE FROM likes WHERE user_id = $1 AND post_id = $2"
	res, err := db.Exec(query, userID, postID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not delete a like: %v", err), http.StatusInternalServerError)
		return
	}

	count, err := res.RowsAffected()
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not delete a like: %v", err), http.StatusInternalServerError)
		return
	}
	if count != 1 {
		http.Error(w, "No row found to delete", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

func CreateMuting(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var body createMutingRequestBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Request body was invalid: %v", err), http.StatusBadRequest)
		return
	}

	sourceUserID := r.PathValue("id")

	query := `INSERT INTO mutes (source_user_id, target_user_id) VALUES ($1, $2)`

	_, err = db.Exec(query, sourceUserID, body.TargetUserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not create muting: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteMuting(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	sourceUserID := r.PathValue("source_user_id")
	targetUserID := r.PathValue("target_user_id")

	query := `DELETE FROM mutes WHERE source_user_id = $1 AND target_user_id = $2`
	res, err := db.Exec(query, sourceUserID, targetUserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not delete muting: %v", err), http.StatusInternalServerError)
		return
	}
	count, err := res.RowsAffected()
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not delete muting: %v", err), http.StatusInternalServerError)
		return
	}
	if count != 1 {
		http.Error(w, "No row found to delete.", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func CreateBlocking(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var body createBlockingRequestBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Request body was invalid: %v", err), http.StatusBadRequest)
		return
	}

	sourceUserID := r.PathValue("id")
	targetUserID := body.TargetUserID

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not start transaction: %v", err), http.StatusInternalServerError)
		return
	}

	defer tx.Rollback()

	query := `INSERT INTO blocks (source_user_id, target_user_id) VALUES ($1, $2)`
	_, err = tx.Exec(query, sourceUserID, targetUserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not create blocking: %v", err), http.StatusInternalServerError)
		return
	}

	query = `DELETE FROM followships WHERE source_user_id = $1 AND target_user_id = $2`
	_, err = tx.Exec(query, sourceUserID, targetUserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not delete followship from source to target: %v", err), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(query, targetUserID, sourceUserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not delete followship from target to source: %v", err), http.StatusInternalServerError)
		return
	}

	query = `DELETE FROM mutes WHERE source_user_id = $1 AND target_user_id = $2`
	_, err = tx.Exec(query, sourceUserID, targetUserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not delete mute: %v", err), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not commit transaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteBlocking(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	sourceUserID := r.PathValue("source_user_id")
	targetUserID := r.PathValue("target_user_id")

	query := `DELETE FROM blocks WHERE source_user_id = $1 AND target_user_id = $2`
	res, err := db.Exec(query, sourceUserID, targetUserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not delete blocking: %v", err), http.StatusInternalServerError)
		return
	}
	count, err := res.RowsAffected()
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not delete blocking: %v", err), http.StatusInternalServerError)
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

func DeleteRepost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	postID := r.PathValue("post_id")
	userID := r.PathValue("user_id")

	query := `DELETE FROM reposts WHERE post_id = $1 AND user_id = $2`

	res, err := db.Exec(query, postID, userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not delete a repost: (user id: %s, post id: %s)\n", userID, postID), http.StatusInternalServerError)
		return
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not delete a repost: (user id: %s, post id: %s)\n", userID, postID), http.StatusInternalServerError)
		return
	}
	if cnt != 1 {
		http.Error(w, fmt.Sprintf("No row found to delete: (user id: %s, post id: %s)\n", userID, postID), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetPosts gets posts from posts table.
func GetPosts(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	query := `SELECT * FROM posts`
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, fmt.Sprintln("Could not get posts"), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []entities.Post
	for rows.Next() {
		var (
			id         uuid.UUID
			user_id    uuid.UUID
			text       string
			created_at time.Time
		)
		if err := rows.Scan(&id, &user_id, &text, &created_at); err != nil {
			http.Error(w, fmt.Sprintln("Could not get posts"), http.StatusInternalServerError)
			return
		}

		post := entities.Post{
			ID:        id,
			UserID:    user_id,
			Text:      text,
			CreatedAt: created_at,
		}
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(posts); err != nil {
		http.Error(w, "Failed to convert to json", http.StatusInternalServerError)
		return
	}
}
