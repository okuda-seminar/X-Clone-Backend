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
	domainerrors "x-clone-backend/internal/app/errors"
	"x-clone-backend/internal/app/usecases"
	"x-clone-backend/internal/domain/entities"

	"github.com/google/uuid"
)

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

// DeletePost deletes a post with the specified post ID.
// If the post doesn't exist, it returns 404 error.
func DeletePost(w http.ResponseWriter, r *http.Request, db *sql.DB, mu *sync.Mutex, usersChan *map[string]chan entities.TimelineEvent) {
	postID := r.PathValue("postID")
	slog.Info(fmt.Sprintf("DELETE /api/posts was called with %s.", postID))

	query := `DELETE FROM posts WHERE id = $1 RETURNING user_id, text, created_at`
	var post entities.Post

	err := db.QueryRow(query, postID).Scan(&post.UserID, &post.Text, &post.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, fmt.Sprintf("No row found to delete (ID: %s)\n", postID), http.StatusNotFound)
			return
		}

		http.Error(w, fmt.Sprintf("Could not delete a post (ID: %s)\n", postID), http.StatusInternalServerError)
		return
	}

	post.ID, err = uuid.Parse(postID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not delete a post (ID: %s)\n", postID), http.StatusInternalServerError)
		return
	}

	go func(userID uuid.UUID, userChan *map[string]chan entities.TimelineEvent) {
		var posts []*entities.Post
		posts = append(posts, &post)
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
				userChan <- entities.TimelineEvent{EventType: entities.PostDeleted, Posts: posts}
			}
			mu.Unlock()
		}

	}(post.UserID, usersChan)

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
func CreateRepost(w http.ResponseWriter, r *http.Request, db *sql.DB, mu *sync.Mutex, usersChan *map[string]chan entities.TimelineEvent) {
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

	var (
		userID    uuid.UUID
		text      string
		createdAt time.Time
	)

	query = `SELECT user_id, text, created_at FROM posts WHERE id=$1`
	err = db.QueryRow(query, body.PostID).Scan(&userID, &text, &createdAt)
	if err != nil {
		http.Error(w, fmt.Sprintln("Could not fetch the original post for repost."), http.StatusInternalServerError)
		return
	}

	post := entities.Post{
		ID:        body.PostID,
		UserID:    userID,
		Text:      text,
		CreatedAt: createdAt,
	}

	go func(userID uuid.UUID, userChan *map[string]chan entities.TimelineEvent) {
		var posts []*entities.Post
		posts = append(posts, &post)
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
				userChan <- entities.TimelineEvent{EventType: entities.RepostCreated, Posts: posts}
			}
			mu.Unlock()
		}

	}(body.UserID, usersChan)

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
func GetReverseChronologicalHomeTimeline(w http.ResponseWriter, r *http.Request, u usecases.GetUserAndFolloweePostsUsecase, mu *sync.Mutex, usersChan *map[string]chan entities.TimelineEvent) {
	userID := r.PathValue("id")
	posts, err := u.GetUserAndFolloweePosts(userID)
	if err != nil {
		http.Error(w, fmt.Sprintln("Could not get posts"), http.StatusInternalServerError)
		return
	}

	mu.Lock()
	if _, exists := (*usersChan)[userID]; !exists {
		(*usersChan)[userID] = make(chan entities.TimelineEvent, 1)
	}
	userChan := (*usersChan)[userID]
	mu.Unlock()

	userChan <- entities.TimelineEvent{EventType: entities.TimelineAccessed, Posts: posts}

	flusher, _ := w.(http.Flusher)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		select {
		case event := <-userChan:
			jsonData, err := json.Marshal(event)
			if err != nil {
				log.Println(err)
				return
			}

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
