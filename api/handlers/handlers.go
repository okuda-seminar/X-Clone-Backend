package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"time"
	"x-clone-backend/api/transfers"
	"x-clone-backend/domain/entities"
	domainerrors "x-clone-backend/domain/errors"
	openapi "x-clone-backend/gen"
	"x-clone-backend/usecases"

	"github.com/google/uuid"
)

// CreateUser creates a new user with the specified useranme and display name,
// then, inserts it into users table.
func CreateUser(w http.ResponseWriter, r *http.Request, u usecases.CreateUserUsecase) {
	var body openapi.CreateUserRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintln("Request body was invalid."), http.StatusBadRequest)
		return
	}

	user, err := u.CreateUser(body.Username, body.DisplayName, body.Password)
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

// DeleteUser deletes a user with the specified user ID.
// If a target user does not exist, it returns 404.
func DeleteUserByID(w http.ResponseWriter, r *http.Request, u usecases.DeleteUserUsecase) {
	userID := r.PathValue("userID")

	slog.Info(fmt.Sprintf("DELETE /api/users was called with %s.", userID))

	err := u.DeleteUser(userID)
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrUserNotFound):
			http.Error(w, fmt.Sprintf("No row found to delete (ID: %s)\n", userID), http.StatusNotFound)
		default:
			http.Error(w, fmt.Sprintf("Could not delete a user (ID: %s)\n", userID), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// FindUserByID finds a user with the specified ID.
func FindUserByID(w http.ResponseWriter, r *http.Request, u usecases.GetSpecificUserUsecase) {
	userID := r.PathValue("userID")

	slog.Info("GET /api/users/{userID} was called.")

	user, err := u.GetSpecificUser(userID)
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
func CreatePost(w http.ResponseWriter, r *http.Request, db *sql.DB, mu *sync.Mutex, usersChan *map[string]chan []byte) {
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

	go func(userID uuid.UUID, userChan *map[string]chan []byte) {
		var posts []entities.Post
		posts = append(posts, post)
		jsonData, err := json.Marshal(posts)
		if err != nil {
			log.Fatalln(err)
			return
		}
		query = `SELECT source_user_id FROM followships WHERE target_user_id=$1`
		rows, err := db.Query(query, userID.String())
		if err != nil {
			log.Fatalln(err)
			return
		}

		var ids []uuid.UUID
		for rows.Next() {
			var id uuid.UUID
			if err := rows.Scan(&id); err != nil {
				log.Fatalln(err)
				return
			}

			ids = append(ids, id)
		}
		ids = append(ids, userID)
		for _, id := range ids {
			mu.Lock()
			if userChan, ok := (*usersChan)[id.String()]; ok {
				userChan <- jsonData
			}
			mu.Unlock()
		}

	}(body.UserID, usersChan)

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
func LikePost(w http.ResponseWriter, r *http.Request, u usecases.LikePostUsecase) {
	var body likePostRequestBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintln("Request body was invalid."), http.StatusBadRequest)
		return
	}

	userID := r.PathValue("id")

	slog.Info(fmt.Sprintf("POST /api/users/{id}/likes was called with %s.", userID))

	err = u.LikePost(userID, body.PostID)
	if err != nil {
		http.Error(w, fmt.Sprintln("Could not create a like."), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UnlikePost(w http.ResponseWriter, r *http.Request, u usecases.UnlikePostUsecase) {
	userID := r.PathValue("id")
	postID := r.PathValue("post_id")

	slog.Info(fmt.Sprintf("DELETE /api/users/{id}/likes/{post_id} was called with %s and %s.", userID, postID))

	err := u.UnlikePost(userID, postID)
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrLikeNotFound):
			http.Error(w, "No row found to delete", http.StatusNotFound)
		default:
			http.Error(w, fmt.Sprintf("Could not delete a like: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func CreateFollowship(w http.ResponseWriter, r *http.Request, u usecases.FollowUserUsecase) {
	var body createFollowshipRequestBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintln("Request body was invalid."), http.StatusBadRequest)
		return
	}

	sourceUserID := r.PathValue("id")

	err = u.FollowUser(sourceUserID, body.TargetUserID)
	if err != nil {
		http.Error(w, fmt.Sprintln("Could not create followship."), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteFollowship(w http.ResponseWriter, r *http.Request, u usecases.UnfollowUserUsecase) {
	sourceUserID := r.PathValue("source_user_id")
	targetUserID := r.PathValue("target_user_id")

	err := u.UnfollowUser(sourceUserID, targetUserID)
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrFollowshipNotFound):
			http.Error(w, "No row found to delete", http.StatusNotFound)
		default:
			http.Error(w, "Could not delete followship.", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func CreateMuting(w http.ResponseWriter, r *http.Request, u usecases.MuteUserUsecase) {
	var body createMutingRequestBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Request body was invalid: %v", err), http.StatusBadRequest)
		return
	}

	sourceUserID := r.PathValue("id")

	err = u.MuteUser(sourceUserID, body.TargetUserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not create muting: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteMuting(w http.ResponseWriter, r *http.Request, u usecases.UnmuteUserUsecase) {
	sourceUserID := r.PathValue("source_user_id")
	targetUserID := r.PathValue("target_user_id")

	err := u.UnmuteUser(sourceUserID, targetUserID)
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrMuteNotFound):
			http.Error(w, "No row found to delete", http.StatusNotFound)
		default:
			http.Error(w, "Could not delete mute.", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func CreateBlocking(w http.ResponseWriter, r *http.Request, u usecases.BlockUserUsecase) {
	var body createBlockingRequestBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Request body was invalid: %v", err), http.StatusBadRequest)
		return
	}

	sourceUserID := r.PathValue("id")
	targetUserID := body.TargetUserID

	err = u.BlockUser(sourceUserID, targetUserID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not create block: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteBlocking(w http.ResponseWriter, r *http.Request, u usecases.UnblockUserUsecase) {
	sourceUserID := r.PathValue("source_user_id")
	targetUserID := r.PathValue("target_user_id")

	err := u.UnblockUser(sourceUserID, targetUserID)
	if err != nil {
		switch {
		case errors.Is(err, domainerrors.ErrBlockNotFound):
			http.Error(w, "No row found to delete", http.StatusNotFound)
		default:
			http.Error(w, "Could not delete blocking.", http.StatusInternalServerError)
		}
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

// GetUserPostsTimeline gets posts by a single user, specified by the requested user ID.
func GetUserPostsTimeline(w http.ResponseWriter, r *http.Request, u usecases.GetSpecificUserPostsUsecase) {
	userID := r.PathValue("id")
	posts, err := u.GetSpecificUserPosts(userID)
	if err != nil {
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(posts); err != nil {
		http.Error(w, "Failed to convert to json", http.StatusInternalServerError)
		return
	}
}

// GetReverseChronologicalHomeTimeline gets posts whose user_id is user or following user from posts table.
func GetReverseChronologicalHomeTimeline(w http.ResponseWriter, r *http.Request, u usecases.GetUserAndFolloweePostsUsecase, mu *sync.Mutex, usersChan *map[string]chan []byte) {
	userID := r.PathValue("id")
	posts, err := u.GetUserAndFolloweePosts(userID)
	if err != nil {
		http.Error(w, fmt.Sprintln("Could not get posts"), http.StatusInternalServerError)
		return
	}

	mu.Lock()
	if _, exists := (*usersChan)[userID]; !exists {
		(*usersChan)[userID] = make(chan []byte, 1)
	}
	userChan := (*usersChan)[userID]
	mu.Unlock()

	jsonData, err := json.Marshal(posts)
	if err != nil {
		log.Println(err)
		return
	}
	userChan <- jsonData

	flusher, _ := w.(http.Flusher)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		select {
		case jsonData := <-userChan:
			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			flusher.Flush()
		case <-r.Context().Done():
			mu.Lock()
			delete(*usersChan, userID)
			mu.Unlock()
			return
		}
	}
}
