package infrastructure

import (
	"database/sql"
	"time"
	"x-clone-backend/domain/entities"
	"x-clone-backend/domain/errors"
	"x-clone-backend/domain/repositories"

	"github.com/google/uuid"
)

type UsersRepository struct {
	DB *sql.DB
}

func NewUsersRepository(db *sql.DB) repositories.UsersRepositoryInterface {
	return &UsersRepository{db}
}

func (r *UsersRepository) CreateUser(username, displayName string) (entities.User, error) {
	query := `INSERT INTO users (username, display_name, bio, is_private) VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at`

	var (
		id                   uuid.UUID
		createdAt, updatedAt time.Time
	)

	err := r.DB.QueryRow(query, username, displayName, "", false).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return entities.User{}, err
	}

	user := entities.User{
		ID:          id,
		Username:    username,
		DisplayName: displayName,
		Bio:         "",
		IsPrivate:   false,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
	return user, nil
}

func (r *UsersRepository) DeleteUser(userID string) error {
	query := `DELETE FROM users WHERE id = $1`
	res, err := r.DB.Exec(query, userID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.ErrUserNotFound
	}

	return nil
}

func (r *UsersRepository) GetSpecificUser(userID string) (entities.User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	row := r.DB.QueryRow(query, userID)

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
	return user, err
}

func (r *UsersRepository) LikePost(userID string, postID uuid.UUID) error {
	query := "INSERT INTO likes (user_id, post_id) VALUES ($1, $2)"

	_, err := r.DB.Exec(query, userID, postID)
	return err
}

func (r *UsersRepository) UnlikePost(userID string, postID string) error {
	query := "DELETE FROM likes WHERE user_id = $1 AND post_id = $2"
	res, err := r.DB.Exec(query, userID, postID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.ErrLikeNotFound
	}

	return nil
}
