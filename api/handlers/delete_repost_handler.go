package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
)

type DeleteRepostHandler struct {
	db *sql.DB
}

func NewDeleteRepostHandler(db *sql.DB) DeleteRepostHandler {
	return DeleteRepostHandler{
		db: db,
	}
}

// DeleteRepost deletes a repost with the specified post ID.
// If the post doesn't exist, it returns 404 error.
func (h *DeleteRepostHandler) DeleteRepost(w http.ResponseWriter, r *http.Request, userID string, postID string) {
	query := `DELETE FROM reposts WHERE post_id = $1 AND user_id = $2`

	res, err := h.db.Exec(query, postID, userID)
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
