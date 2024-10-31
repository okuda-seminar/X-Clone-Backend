package usecases

import (
	"database/sql"
	"x-clone-backend/domain/repositories"
)

type UnmuteUserUsecase interface {
	UnmuteUser(sourceUserID, targetUserID string) error
}

type unmuteUserUsecase struct {
	usersRepository repositories.UsersRepositoryInterface
}

func NewUnmuteUserUsecase(usersRepository repositories.UsersRepositoryInterface) UnmuteUserUsecase {
	return &unmuteUserUsecase{usersRepository: usersRepository}
}

func (p *unmuteUserUsecase) UnmuteUser(sourceUserID, targetUserID string) error {
	err := p.usersRepository.WithTransaction(func(tx *sql.Tx) error {
		return p.usersRepository.UnmuteUser(tx, sourceUserID, targetUserID)
	})
	return err
}
