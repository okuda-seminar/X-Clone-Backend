package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"x-clone-backend/entities"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

func (s *HandlersTestSuite) TestCreateUser() {
	tests := []struct {
		name         string
		body         string
		expectedCode int
	}{
		{
			name:         "create user",
			body:         `{ "username": "test", "display_name": "test" }`,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid JSON body",
			body:         `{ "username": "` + "test",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid body",
			body:         `{ "invalid": "test" }`,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "duplicated username",
			body:         `{ "username": "test", "display_name": "duplicated" }`,
			expectedCode: http.StatusConflict,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("POST", "/api/users", strings.NewReader(test.body))
		rr := httptest.NewRecorder()

		CreateUser(rr, req, s.db)

		if rr.Code != test.expectedCode {
			s.T().Errorf("%s: wrong code returned; expected %d, but got %d", test.name, test.expectedCode, rr.Code)
		}
	}
}

func (s *HandlersTestSuite) TestLikePost() {
	// LikePost must use existing user ID and post ID
	// from the users and posts table.
	// Therefore, users and posts are created
	// for testing purposes to obtain these IDs.
	authorUserId := s.getTestUserId(`{ "username": "author", "display_name": "author" }`)
	likerUserId := s.getTestUserId(`{ "username": "liker", "display_name": "liker" }`)
	postId := s.getTestPostId(fmt.Sprintf(`{ "user_id": "%s", "text": "test post"}`, authorUserId))

	tests := []struct {
		name         string
		userId       string
		body         string
		expectedCode int
	}{
		{
			name:         "create like for a own post",
			userId:       authorUserId,
			body:         fmt.Sprintf(`{ "post_id": "%s" }`, postId),
			expectedCode: http.StatusCreated,
		},
		{
			name:         "create like for other's post",
			userId:       likerUserId,
			body:         fmt.Sprintf(`{ "post_id": "%s" }`, postId),
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid JSON body",
			userId:       likerUserId,
			body:         fmt.Sprintf(`{ "post_id": "%s"`, postId),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid body",
			userId:       likerUserId,
			body:         `{ "invalid": "test" }`,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "non-existent user id",
			userId:       uuid.New().String(),
			body:         fmt.Sprintf(`{ "post_id": "%s" }`, postId),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "non-existent post id",
			userId:       likerUserId,
			body:         fmt.Sprintf(`{ "post_id": "%s" }`, uuid.New().String()),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "duplicated like",
			userId:       likerUserId,
			body:         fmt.Sprintf(`{ "post_id": "%s" }`, postId),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest(
			"POST",
			"/api/users/{id}/likes",
			strings.NewReader(test.body),
		)
		req.SetPathValue("id", test.userId)

		rr := httptest.NewRecorder()
		LikePost(rr, req, s.db)

		if rr.Code != test.expectedCode {
			s.T().Errorf(
				"%s: wrong code returned; expected %d, but got %d",
				test.name,
				test.expectedCode,
				rr.Code,
			)
		}
	}
}

func (s *HandlersTestSuite) TestCreateMuting() {
	// CreateMuting must use existing user IDs from the user table
	// for both the source user ID and target user ID.
	// Therefore, users are created for testing purposes to obtain these IDs.
	sourceUserId := s.getTestUserId(`{ "username": "test", "display_name": "test" }`)
	targetUserId := s.getTestUserId(`{ "username": "test2", "display_name": "test2" }`)

	tests := []struct {
		name         string
		body         string
		expectedCode int
	}{
		{
			name:         "create muting",
			body:         `{ "target_user_id": "` + targetUserId + `" }`,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid JSON body",
			body:         `{ "target_user_id": "` + targetUserId,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid body",
			body:         `{ "invalid": "test" }`,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "duplicated muting",
			body:         `{ "target_user_id": "` + targetUserId + `" }`,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest(
			"POST",
			"/api/users/{id}/muting",
			strings.NewReader(test.body),
		)
		req.SetPathValue("id", sourceUserId)

		rr := httptest.NewRecorder()
		CreateMuting(rr, req, s.db)

		if rr.Code != test.expectedCode {
			s.T().Errorf(
				"%s: wrong code returned; expected %d, but got %d",
				test.name,
				test.expectedCode,
				rr.Code,
			)
		}
	}
}

func (s *HandlersTestSuite) TestCreateRepost() {
	userId := s.getTestUserId(`{ "username": "test", "display_name": "test" }`)
	postId := s.getTestPostId(fmt.Sprintf(`{ "user_id": "%s", "text": "test" }`, userId))

	tests := []struct {
		name         string
		body         string
		expectedCode int
	}{
		{
			name:         "create repost",
			body:         fmt.Sprintf(`{ "post_id": "%s", "user_id": "%s" }`, postId, userId),
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid JSON body",
			body:         fmt.Sprintf(`{ "post_id": "%s", "user_id": "%s }`, postId, userId),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid body",
			body:         fmt.Sprintf(`{ "post_id": "%s" }`, postId),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "non-existent user id",
			body:         fmt.Sprintf(`{ "post_id": "%s", "user_id": "%s" }`, postId, uuid.New()),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "non-existent post id",
			body:         fmt.Sprintf(`{ "post_id": "%s", "user_id": "%s" }`, uuid.New(), userId),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "duplicated repost",
			body:         fmt.Sprintf(`{ "post_id": "%s", "user_id": "%s" }`, postId, userId),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest(
			"POST",
			"/api/posts/reposts",
			strings.NewReader(test.body),
		)
		rr := httptest.NewRecorder()
		CreateRepost(rr, req, s.db)

		if rr.Code != test.expectedCode {
			s.T().Errorf("%s: wrong code returned; expected %d, but got %d", test.name, test.expectedCode, rr.Code)
		}
	}
}

func (s *HandlersTestSuite) getTestUserId(body string) string {
	req := httptest.NewRequest(
		"POST",
		"/api/users",
		strings.NewReader(body),
	)
	rr := httptest.NewRecorder()
	CreateUser(rr, req, s.db)

	var user entities.User
	_ = json.NewDecoder(rr.Body).Decode(&user)
	sourceUserId := user.ID.String()
	return sourceUserId
}

func (s *HandlersTestSuite) getTestPostId(body string) string {
	req := httptest.NewRequest(
		"POST",
		"/api/posts",
		strings.NewReader(body),
	)
	rr := httptest.NewRecorder()
	CreatePost(rr, req, s.db)

	var post entities.Post
	_ = json.NewDecoder(rr.Body).Decode(&post)
	postId := post.ID.String()
	return postId
}

// TestHandlersTestSuite runs all of the tests attached to HandlersTestSuite.
func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}
