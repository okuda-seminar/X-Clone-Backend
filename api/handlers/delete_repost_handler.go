package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	"x-clone-backend/internal/domain/entities"

	"github.com/google/uuid"
)

type DeleteRepostHandler struct {
	db        *sql.DB
	mu        *sync.Mutex
	usersChan *map[string]chan entities.TimelineEvent
}

func NewDeleteRepostHandler(db *sql.DB, mu *sync.Mutex, usersChan *map[string]chan entities.TimelineEvent) DeleteRepostHandler {
	return DeleteRepostHandler{
		db:        db,
		mu:        mu,
		usersChan: usersChan,
	}
}

// DeleteRepost deletes a repost with the specified post ID.
// If the post doesn't exist, it returns 404 error.
func (h *DeleteRepostHandler) DeleteRepost(w http.ResponseWriter, r *http.Request, userIDStr string, parentIDStr string) {
	var body deleteRepostRequestBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintln("Request body was invalid."), http.StatusBadRequest)
		return
	}

	query := `DELETE FROM reposts WHERE id = $1 RETURNING text, created_at`

	var (
		text      string
		createdAt time.Time
	)

	err = h.db.QueryRow(query, body.RepostID).Scan(&text, &createdAt)
	if err == sql.ErrNoRows {
		http.Error(w, fmt.Sprintf("No row found to delete: (repost id: %s)\n", body.RepostID), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not delete a repost: (repost id: %s)\n", body.RepostID), http.StatusInternalServerError)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not parse a userID (ID: %s)\n", userIDStr), http.StatusBadRequest)
		return
	}

	parentID, err := uuid.Parse(parentIDStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not parse a parentID (ID: %s)\n", parentIDStr), http.StatusBadRequest)
		return
	}

	repost := entities.Repost{
		ID:        body.RepostID,
		ParentID:  parentID,
		UserID:    userID,
		Text:      text,
		CreatedAt: createdAt,
	}

	go func(userID uuid.UUID, usersChan *map[string]chan entities.TimelineEvent) {
		var reposts []*entities.Repost
		reposts = append(reposts, &repost)
		query = `SELECT source_user_id FROM followships WHERE target_user_id=$1`
		rows, err := h.db.Query(query, userID.String())
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
			h.mu.Lock()
			if userChan, ok := (*usersChan)[id.String()]; ok {
				userChan <- entities.TimelineEvent{EventType: entities.RepostDeleted, Reposts: reposts}
			}
			h.mu.Unlock()
		}
	}(userID, h.usersChan)

	w.WriteHeader(http.StatusNoContent)
}
