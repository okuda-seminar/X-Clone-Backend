package usecases

import (
	"x-clone-backend/internal/app/errors"
	"x-clone-backend/internal/domain/repositories"
)

type BlockUserUsecase interface {
	BlockUser(sourceUserID, targetUserID string) error
}

type blockUserUsecase struct {
	usersRepository repositories.UsersRepositoryInterface
}

func NewBlockUserUsecase(usersRepository repositories.UsersRepositoryInterface) BlockUserUsecase {
	return &blockUserUsecase{usersRepository: usersRepository}
}

func (p *blockUserUsecase) BlockUser(sourceUserID, targetUserID string) error {
	if err := p.usersRepository.BlockUser(sourceUserID, targetUserID); err != nil {
		return err
	}
	if err := p.usersRepository.UnfollowUser(sourceUserID, targetUserID); err != nil {
		if err != errors.ErrFollowshipNotFound {
			return err
		}
	}
	if err := p.usersRepository.UnfollowUser(targetUserID, sourceUserID); err != nil {
		if err != errors.ErrFollowshipNotFound {
			return err
		}
	}
	if err := p.usersRepository.UnmuteUser(sourceUserID, targetUserID); err != nil {
		if err != errors.ErrMuteNotFound {
			return err
		}
	}
	return nil
}
