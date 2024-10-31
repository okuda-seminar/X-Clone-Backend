package usecases

import (
	"database/sql"
	"x-clone-backend/domain/repositories"
)

type FollowUserUsecase interface {
	FollowUser(sourceUserID, targetUserID string) error
}

type followUserUsecase struct {
	usersRepository repositories.UsersRepositoryInterface
}

func NewFollowUserUsecase(usersRepository repositories.UsersRepositoryInterface) FollowUserUsecase {
	return &followUserUsecase{usersRepository: usersRepository}
}

func (p *followUserUsecase) FollowUser(sourceUserID, targetUserID string) error {
	err := p.usersRepository.WithTransaction(func(tx *sql.Tx) error {
		return p.usersRepository.FollowUser(tx, sourceUserID, targetUserID)
	})
	return err
}
