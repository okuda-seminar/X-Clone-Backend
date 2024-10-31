package usecases

import (
	"database/sql"
	"x-clone-backend/domain/repositories"
)

type MuteUserUsecase interface {
	MuteUser(sourceUserID, targetUserID string) error
}

type muteUserUsecase struct {
	usersRepository repositories.UsersRepositoryInterface
}

func NewMuteUserUsecase(usersRepository repositories.UsersRepositoryInterface) MuteUserUsecase {
	return &muteUserUsecase{usersRepository: usersRepository}
}

func (p *muteUserUsecase) MuteUser(sourceUserID, targetUserID string) error {
	err := p.usersRepository.WithTransaction(func(tx *sql.Tx) error {
		return p.usersRepository.MuteUser(tx, sourceUserID, targetUserID)
	})
	return err
}
