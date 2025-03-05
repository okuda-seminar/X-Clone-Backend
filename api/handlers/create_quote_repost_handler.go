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

type CreateQuoteRepostHandler struct {
	db        *sql.DB
	mu        *sync.Mutex
	usersChan *map[string]chan entities.TimelineEvent
}

func NewCreateQuoteRepostHandler(db *sql.DB, mu *sync.Mutex, usersChan *map[string]chan entities.TimelineEvent) CreateQuoteRepostHandler {
	return CreateQuoteRepostHandler{
		db:        db,
		mu:        mu,
		usersChan: usersChan,
	}
}

// CreateQuoteRepost creates a new quote repost with the specified post_id and user_id,
// then, inserts it into reposts table.
func (h *CreateQuoteRepostHandler) CreateQuoteRepost(w http.ResponseWriter, r *http.Request, userIDStr string) {
	var body createQuoteRepostRequestBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		http.Error(w, fmt.Sprintln("Request body was invalid."), http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not parse a userID (ID: %s)\n", userIDStr), http.StatusBadRequest)
		return
	}

	query := `
		SELECT 
			r.id IS NOT NULL AS is_parent_repost
		FROM users u
		LEFT JOIN reposts r ON r.id = $2
		WHERE u.id = $1
	`
	var isParentRepost bool
	err = h.db.QueryRow(query, userID, body.PostID).Scan(&isParentRepost)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found.", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Database query error.", http.StatusInternalServerError)
		return
	}

	if isParentRepost {
		query = `INSERT INTO reposts (user_id, parent_repost_id, is_quote, text) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	} else {
		query = `INSERT INTO reposts (user_id, parent_post_id, is_quote, text) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	}

	var (
		id        uuid.UUID
		createdAt time.Time
	)

	isQuote := true

	err = h.db.QueryRow(query, userID, body.PostID, isQuote, body.Text).Scan(&id, &createdAt)
	if err != nil {
		http.Error(w, fmt.Sprintln("Could not create a quote repost."), http.StatusInternalServerError)
		return
	}

	quoteRepost := entities.Repost{
		ID:        id,
		ParentID:  body.PostID,
		UserID:    userID,
		Text:      body.Text,
		CreatedAt: createdAt,
	}

	go func(userID uuid.UUID, userChan *map[string]chan entities.TimelineEvent) {
		var quoteReposts []*entities.Repost
		quoteReposts = append(quoteReposts, &quoteRepost)
		query := `SELECT source_user_id FROM followships WHERE target_user_id = $1`
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
			if userChan, ok := (*h.usersChan)[id.String()]; ok {
				userChan <- entities.TimelineEvent{EventType: entities.QuoteRepostCreated, Reposts: quoteReposts}
			}
			h.mu.Unlock()
		}

	}(userID, h.usersChan)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)
	err = encoder.Encode(&quoteRepost)
	if err != nil {
		http.Error(w, fmt.Sprintln("Could not encode response."), http.StatusInternalServerError)
		return
	}
}
