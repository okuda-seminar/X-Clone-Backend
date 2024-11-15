package usecases

import (
	"fmt"
	"x-clone-backend/internal/app/services"
	"x-clone-backend/internal/domain/entities"
	"x-clone-backend/internal/domain/repositories"
)

type CreateUserUsecase interface {
	CreateUser(username, displayName, password string) (entities.User, error)
}

type createUserUsecase struct {
	usersRepository repositories.UsersRepositoryInterface
}

func NewCreateUserUsecase(usersRepository repositories.UsersRepositoryInterface) CreateUserUsecase {
	return &createUserUsecase{usersRepository: usersRepository}
}

func (p *createUserUsecase) CreateUser(username, displayName, password string) (entities.User, error) {
	var user entities.User
	hashedPassword, err := services.HashPassword(password)
	if err != nil {
		return entities.User{}, fmt.Errorf("could not hash the password: %w", err)
	}
	user, err = p.usersRepository.CreateUser(nil, username, displayName, hashedPassword)
	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}
