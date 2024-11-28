package infrastructure

import (
	"x-clone-backend/internal/app/errors"
	"x-clone-backend/internal/domain/entities"
	"x-clone-backend/internal/domain/repositories"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

type UsersRepository struct {
	DB *gorm.DB
}

func NewUsersRepository(db *gorm.DB) repositories.UsersRepositoryInterface {
	return &UsersRepository{db}
}

func (r *UsersRepository) CreateUser(username, displayName, password string) (entities.User, error) {
	user := entities.User{
		Username:    username,
		DisplayName: displayName,
		Bio:         "",
		IsPrivate:   false,
		Password:    password,
	}

	if err := r.DB.Create(&user).Error; err != nil {
		return entities.User{}, err
	}

	return user, nil
}

func (r *UsersRepository) DeleteUser(userID string) error {
	var user entities.User
	res := r.DB.Delete(&user, "id = ?", userID)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected != 1 {
		return errors.ErrUserNotFound
	}

	return nil
}

func (r *UsersRepository) GetSpecificUser(userID string) (entities.User, error) {
	var user entities.User
	if err := r.DB.First(&user, "id = ?", userID).Error; err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func (r *UsersRepository) LikePost(userID string, postID uuid.UUID) error {
	like := entities.Like{
		UserID: userID,
		PostID: postID,
	}

	if err := r.DB.Create(&like).Error; err != nil {
		return err
	}

	return nil
}

func (r *UsersRepository) UnlikePost(userID string, postID string) error {
	res := r.DB.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&entities.Like{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 1 {
		return errors.ErrLikeNotFound
	}

	return nil
}

func (r *UsersRepository) FollowUser(sourceUserID, targetUserID string) error {
	followship := entities.Followship{
		SourceUserID: sourceUserID,
		TargetUserID: targetUserID,
	}

	if err := r.DB.Create(&followship).Error; err != nil {
		return err
	}

	return nil
}

func (r *UsersRepository) UnfollowUser(sourceUserID, targetUserID string) error {
	res := r.DB.Where("source_user_id = ? AND target_user_id = ?", sourceUserID, targetUserID).Delete(&entities.Followship{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 1 {
		return errors.ErrFollowshipNotFound
	}

	return nil
}

func (r *UsersRepository) MuteUser(sourceUserID, targetUserID string) error {
	mute := entities.Mute{
		SourceUserID: sourceUserID,
		TargetUserID: targetUserID,
	}

	if err := r.DB.Create(&mute).Error; err != nil {
		return err
	}

	return nil
}

func (r *UsersRepository) UnmuteUser(sourceUserID, targetUserID string) error {
	res := r.DB.Where("source_user_id = ? AND target_user_id = ?", sourceUserID, targetUserID).Delete(&entities.Mute{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 1 {
		return errors.ErrMuteNotFound
	}

	return nil
}

func (r *UsersRepository) BlockUser(sourceUserID, targetUserID string) error {
	block := entities.Block{
		SourceUserID: sourceUserID,
		TargetUserID: targetUserID,
	}
	if err := r.DB.Create(&block).Error; err != nil {
		return err
	}

	return nil
}

func (r *UsersRepository) UnblockUser(sourceUserID, targetUserID string) error {
	res := r.DB.Where("source_user_id = ? AND target_user_id = ?", sourceUserID, targetUserID).Delete(&entities.Block{})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected != 1 {
		return errors.ErrBlockNotFound
	}

	return nil
}
