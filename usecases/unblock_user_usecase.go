package usecases

import (
	"database/sql"
	"x-clone-backend/domain/repositories"
)

type UnblockUserUsecase interface {
	UnblockUser(sourceUserID, targetUserID string) error
}

type unblockUserUsecase struct {
	usersRepository repositories.UsersRepositoryInterface
}

func NewUnblockUserUsecase(usersRepository repositories.UsersRepositoryInterface) UnblockUserUsecase {
	return &unblockUserUsecase{usersRepository: usersRepository}
}

func (p *unblockUserUsecase) UnblockUser(sourceUserID, targetUserID string) error {
	err := p.usersRepository.WithTransaction(func(tx *sql.Tx) error {
		return p.usersRepository.UnblockUser(tx, sourceUserID, targetUserID)
	})
	return err
}
