package usecases

import (
	"database/sql"
	"x-clone-backend/domain/repositories"
)

type DeleteUserUsecase interface {
	DeleteUser(userID string) error
}

type deleteUserUsecase struct {
	usersRepository repositories.UsersRepositoryInterface
}

func NewDeleteUserUsecase(usersRepository repositories.UsersRepositoryInterface) DeleteUserUsecase {
	return &deleteUserUsecase{usersRepository: usersRepository}
}

func (p *deleteUserUsecase) DeleteUser(userID string) error {
	err := p.usersRepository.WithTransaction(func(tx *sql.Tx) error {
		return p.usersRepository.DeleteUser(tx, userID)
	})
	return err
}
