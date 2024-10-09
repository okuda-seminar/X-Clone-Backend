package transfers

import (
	"x-clone-backend/domain/entities"
	openapi "x-clone-backend/gen"
)

func ToCreateUserResponse(in *entities.User) *openapi.CreateUserResponse {
	return &openapi.CreateUserResponse{
		Bio:         in.Bio,
		CreatedAt:   in.CreatedAt,
		DisplayName: in.DisplayName,
		Id:          in.ID.String(),
		IsPrivate:   in.IsPrivate,
		UpdatedAt:   in.UpdatedAt,
		Username:    in.Username,
	}
}
