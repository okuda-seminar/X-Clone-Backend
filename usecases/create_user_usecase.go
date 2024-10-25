package usecases

import (
	"x-clone-backend/domain/entities"
	"x-clone-backend/domain/repositories"
)

type CreateUserUsecase interface {
	CreateUser(username, displayName string) (entities.User, error)
}

type createUserUsecase struct {
	usersRepository repositories.UsersRepositoryInterface
}

func NewCreateUserUsecase(usersRepository repositories.UsersRepositoryInterface) CreateUserUsecase {
	return &createUserUsecase{usersRepository: usersRepository}
}

func (p *createUserUsecase) CreateUser(username, displayName string) (entities.User, error) {
	user, err := p.usersRepository.CreateUser(username, displayName)
	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}
