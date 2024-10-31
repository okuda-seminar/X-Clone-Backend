package usecases

import (
	"database/sql"
	"x-clone-backend/domain/repositories"
)

type UnfollowUserUsecase interface {
	UnfollowUser(sourceUserID, targetUserID string) error
}

type unfollowUserUsecase struct {
	usersRepository repositories.UsersRepositoryInterface
}

func NewUnfollowUserUsecase(usersRepository repositories.UsersRepositoryInterface) UnfollowUserUsecase {
	return &unfollowUserUsecase{usersRepository: usersRepository}
}

func (p *unfollowUserUsecase) UnfollowUser(sourceUserID, targetUserID string) error {
	err := p.usersRepository.WithTransaction(func(tx *sql.Tx) error {
		return p.usersRepository.UnfollowUser(tx, sourceUserID, targetUserID)
	})
	return err
}
