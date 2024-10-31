package usecases

import (
	"database/sql"
	"x-clone-backend/domain/repositories"
)

type UnlikePostUsecase interface {
	UnlikePost(userID string, postID string) error
}

type unlikePostUsecase struct {
	usersRepository repositories.UsersRepositoryInterface
}

func NewUnlikePostUsecase(usersRepository repositories.UsersRepositoryInterface) UnlikePostUsecase {
	return &unlikePostUsecase{usersRepository: usersRepository}
}

func (p *unlikePostUsecase) UnlikePost(userID string, postID string) error {
	err := p.usersRepository.WithTransaction(func(tx *sql.Tx) error {
		return p.usersRepository.UnlikePost(tx, userID, postID)
	})
	return err
}
