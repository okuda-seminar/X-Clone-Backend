package repositories

import (
	"x-clone-backend/internal/domain/entities"

	"github.com/google/uuid"
)

type UsersRepositoryInterface interface {
	CreateUser(username, displayName, password string) (entities.User, error)
	DeleteUser(userID string) error
	GetSpecificUser(userID string) (entities.User, error)
	LikePost(userID string, postID uuid.UUID) error
	UnlikePost(userID string, postID string) error
	FollowUser(sourceUserID, targetUserID string) error
	UnfollowUser(sourceUserID, targetUserID string) error
	MuteUser(sourceUserID, targetUserID string) error
	UnmuteUser(sourceUserID, targetUserID string) error
	BlockUser(sourceUserID, targetUserID string) error
	UnblockUser(sourceUserID, targetUserID string) error
}
