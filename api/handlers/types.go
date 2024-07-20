package handlers

import (
	"github.com/google/uuid"
)

// createUserRequestBody is the type of the "CreateUser"
// endpoint request body.
type createUserRequestBody struct {
	Username    string `json:"username,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}

// createPostRequestBody is the type of the "CreatePost"
// endpoint request body.
type createPostRequestBody struct {
	UserID uuid.UUID `json:"user_id,omitempty"`
	Text   string    `json:"text"`
}

// createFollowshipRequestBody is the type of the "CreateFollowship"
// endpoint request body.
type createFollowshipRequestBody struct {
	TargetUserID string `json:"target_user_id"`
}
