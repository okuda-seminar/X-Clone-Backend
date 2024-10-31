package usecases

import (
	"database/sql"
	"x-clone-backend/domain/entities"
	"x-clone-backend/domain/repositories"
)

type GetSpecificUserUsecase interface {
	GetSpecificUser(userID string) (entities.User, error)
}

type getSpecificUserUsecase struct {
	usersRepository repositories.UsersRepositoryInterface
}

func NewGetSpecificUserUsecase(usersRepository repositories.UsersRepositoryInterface) GetSpecificUserUsecase {
	return &getSpecificUserUsecase{usersRepository: usersRepository}
}

func (p *getSpecificUserUsecase) GetSpecificUser(userID string) (entities.User, error) {
	var user entities.User
	err := p.usersRepository.WithTransaction(func(tx *sql.Tx) error {
		var err error
		user, err = p.usersRepository.GetSpecificUser(tx, userID)
		return err
	})

	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}
