package usecases

import (
	"database/sql"
	"x-clone-backend/domain/repositories"

	"github.com/google/uuid"
)

type LikePostUsecase interface {
	LikePost(userID string, postID uuid.UUID) error
}

type likePostUsecase struct {
	usersRepository repositories.UsersRepositoryInterface
}

func NewLikePostUsecase(usersRepository repositories.UsersRepositoryInterface) LikePostUsecase {
	return &likePostUsecase{usersRepository: usersRepository}
}

func (p *likePostUsecase) LikePost(userID string, postID uuid.UUID) error {
	err := p.usersRepository.WithTransaction(func(tx *sql.Tx) error {
		return p.usersRepository.LikePost(tx, userID, postID)
	})
	return err
}
