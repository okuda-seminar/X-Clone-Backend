// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package openapi

import (
	"time"
)

// CreatePostRequest defines model for create_post_request.
type CreatePostRequest struct {
	Text   string `json:"text"`
	UserId string `json:"user_id"`
}

// CreatePostResponse defines model for create_post_response.
type CreatePostResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Id        string    `json:"id"`
	Text      string    `json:"text"`
	UserId    string    `json:"user_id"`
}

// CreateQuoteRepostRequest defines model for create_quote_repost_request.
type CreateQuoteRepostRequest struct {
	PostId string `json:"post_id"`
	Text   string `json:"text"`
}

// CreateQuoteRepostResponse defines model for create_quote_repost_response.
type CreateQuoteRepostResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Id        string    `json:"id"`
	ParentId  string    `json:"parent_id"`
	Text      string    `json:"text"`
	UserId    string    `json:"user_id"`
}

// CreateRepostRequest defines model for create_repost_request.
type CreateRepostRequest struct {
	PostId string `json:"post_id"`
}

// CreateRepostResponse defines model for create_repost_response.
type CreateRepostResponse struct {
	CreatedAt time.Time `json:"created_at"`
	Id        string    `json:"id"`
	ParentId  string    `json:"parent_id"`
	Text      string    `json:"text"`
	UserId    string    `json:"user_id"`
}

// CreateUserRequest defines model for create_user_request.
type CreateUserRequest struct {
	DisplayName string `json:"display_name"`

	// Password Password must be between 8 and 15 characters.
	Password string `json:"password"`
	Username string `json:"username"`
}

// CreateUserResponse defines model for create_user_response.
type CreateUserResponse struct {
	Token string `json:"token"`
	User  struct {
		Bio         string    `json:"bio"`
		CreatedAt   time.Time `json:"created_at"`
		DisplayName string    `json:"display_name"`
		Id          string    `json:"id"`
		IsPrivate   bool      `json:"is_private"`
		UpdatedAt   time.Time `json:"updated_at"`
		Username    string    `json:"username"`
	} `json:"user"`
}

// DeleteRepostRequest defines model for delete_repost_request.
type DeleteRepostRequest struct {
	RepostId string `json:"repost_id"`
}

// FindUserByIdResponse defines model for find_user_by_id_response.
type FindUserByIdResponse struct {
	Bio         string    `json:"bio"`
	CreatedAt   time.Time `json:"created_at"`
	DisplayName string    `json:"display_name"`
	Id          string    `json:"id"`
	IsPrivate   bool      `json:"is_private"`
	UpdatedAt   time.Time `json:"updated_at"`
	Username    string    `json:"username"`
}

// GetReverseChronologicalHomeTimelineResponse defines model for get_reverse_chronological_home_timeline_response.
type GetReverseChronologicalHomeTimelineResponse struct {
	Data *struct {
		EventType string `json:"event_type"`
		Posts     struct {
			CreatedAt time.Time `json:"created_at"`
			Id        string    `json:"id"`
			Text      string    `json:"text"`
			UserId    string    `json:"user_id"`
		} `json:"posts"`
		Reposts struct {
			CreatedAt time.Time `json:"created_at"`
			Id        string    `json:"id"`
			ParentId  string    `json:"parent_id"`
			Text      string    `json:"text"`
			UserId    string    `json:"user_id"`
		} `json:"reposts"`
	} `json:"data,omitempty"`
}

// GetUserPostsTimelineResponse defines model for get_user_posts_timeline_response.
type GetUserPostsTimelineResponse = []struct {
	CreatedAt time.Time `json:"created_at"`
	Id        string    `json:"id"`
	Text      string    `json:"text"`
	UserId    string    `json:"user_id"`
}

// CreatePostJSONRequestBody defines body for CreatePost for application/json ContentType.
type CreatePostJSONRequestBody = CreatePostRequest

// CreateUserJSONRequestBody defines body for CreateUser for application/json ContentType.
type CreateUserJSONRequestBody = CreateUserRequest

// CreateQuoteRepostJSONRequestBody defines body for CreateQuoteRepost for application/json ContentType.
type CreateQuoteRepostJSONRequestBody = CreateQuoteRepostRequest

// CreateRepostJSONRequestBody defines body for CreateRepost for application/json ContentType.
type CreateRepostJSONRequestBody = CreateRepostRequest

// DeleteRepostJSONRequestBody defines body for DeleteRepost for application/json ContentType.
type DeleteRepostJSONRequestBody = DeleteRepostRequest
