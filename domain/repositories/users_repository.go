package repositories

import (
	"x-clone-backend/domain/entities"

	"github.com/google/uuid"
)

type UsersRepositoryInterface interface {
	CreateUser(username, displayName string) (entities.User, error)
	DeleteUser(userID string) error
	GetSpecificUser(userID string) (entities.User, error)
	LikePost(userID string, postID uuid.UUID) error
	UnlikePost(userID string, postID string) error
}
