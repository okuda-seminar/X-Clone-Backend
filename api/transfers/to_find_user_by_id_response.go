package transfers

import (
	openapi "x-clone-backend/gen"
	"x-clone-backend/internal/domain/entities"
)

func ToFindUserByIDResponse(in *entities.User) *openapi.FindUserByIdResponse {
	return &openapi.FindUserByIdResponse{
		Bio:         in.Bio,
		CreatedAt:   in.CreatedAt,
		DisplayName: in.DisplayName,
		Id:          in.ID.String(),
		IsPrivate:   in.IsPrivate,
		UpdatedAt:   in.UpdatedAt,
		Username:    in.Username,
	}
}
