package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"x-clone-backend/domain/entities"

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

func (s *HandlersTestSuite) TestDeletePost() {
	userID := s.newTestUser(`{ "username": "test user", "display_name": "test user" }`)
	postID := s.newTestPost(fmt.Sprintf(`{ "user_id": "%s", "text": "test post"}`, userID))

	tests := []struct {
		name         string
		postID       string
		expectedCode int
	}{
		{
			name:         "delete a post successfully with a proper post ID.",
			postID:       postID,
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "fail to delete a post that was already deleted .",
			postID:       postID,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "fail to delete a post with a non-existent post ID.",
			postID:       uuid.New().String(),
			expectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest("DELETE", "/api/posts{postID}",
			nil)
		req.SetPathValue("postID", test.postID)

		rr := httptest.NewRecorder()
		DeletePost(rr, req, s.db)

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

func (s *HandlersTestSuite) TestLikePost() {
	// LikePost must use existing user ID and post ID
	// from the users and posts table.
	// Therefore, users and posts are created
	// for testing purposes to obtain these IDs.
	authorUserID := s.newTestUser(`{ "username": "author", "display_name": "author" }`)
	likerUserID := s.newTestUser(`{ "username": "liker", "display_name": "liker" }`)
	postID := s.newTestPost(fmt.Sprintf(`{ "user_id": "%s", "text": "test post"}`, authorUserID))

	tests := []struct {
		name         string
		userID       string
		body         string
		expectedCode int
	}{
		{
			name:         "like an own post successfully with a proper pair of User and Post",
			userID:       authorUserID,
			body:         fmt.Sprintf(`{ "post_id": "%s" }`, postID),
			expectedCode: http.StatusCreated,
		},
		{
			name:         "like another user's post successfully with a proper pair of User and Post",
			userID:       likerUserID,
			body:         fmt.Sprintf(`{ "post_id": "%s" }`, postID),
			expectedCode: http.StatusCreated,
		},
		{
			name:         "fail to like another user's post with a invalid JSON body",
			userID:       likerUserID,
			body:         fmt.Sprintf(`{ "post_id": "%s"`, postID),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "fail to like a post with a invalid JSON field",
			userID:       likerUserID,
			body:         `{ "invalid": "test" }`,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "fail to like a post with a pair of non-existent User and proper Post",
			userID:       uuid.New().String(),
			body:         fmt.Sprintf(`{ "post_id": "%s" }`, postID),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "fail to like a post with a pair of proper User and non-existent Post",
			userID:       likerUserID,
			body:         fmt.Sprintf(`{ "post_id": "%s" }`, uuid.New().String()),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "fail to like another user's post duplicately with a proper pair of User and Post",
			userID:       likerUserID,
			body:         fmt.Sprintf(`{ "post_id": "%s" }`, postID),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest(
			"POST",
			"/api/users/{id}/likes",
			strings.NewReader(test.body),
		)
		req.SetPathValue("id", test.userID)

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

func (s *HandlersTestSuite) TestUnlikePost() {
	// UnlikePost must use existing user ID and post ID
	// from the users and posts table.
	// Therefore, users and posts are created
	// for testing purposes to obtain these IDs.
	userID := s.newTestUser(`{ "username": "user", "display_name": "user" }`)
	postID := s.newTestPost(fmt.Sprintf(`{ "user_id": "%s", "text": "test post" }`, userID))
	s.newTestLike(userID, postID)

	tests := []struct {
		name         string
		userID       string
		postID       string
		expectedCode int
	}{
		{
			name:         "unlike a post successfully with a proper pair of User and Post.",
			userID:       userID,
			postID:       postID,
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "fail to unlike a post with a pair of non-existent User and proper Post.",
			userID:       uuid.New().String(),
			postID:       postID,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "fail to unlike a post with a pair of proper User and non-existent Post.",
			userID:       userID,
			postID:       uuid.New().String(),
			expectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest(
			"DELETE",
			"/api/users/{id}/likes/{post_id}",
			strings.NewReader(""),
		)
		req.SetPathValue("id", test.userID)
		req.SetPathValue("post_id", test.postID)

		rr := httptest.NewRecorder()
		UnlikePost(rr, req, s.db)

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
	sourceUserID := s.newTestUser(`{ "username": "test", "display_name": "test" }`)
	targetUserID := s.newTestUser(`{ "username": "test2", "display_name": "test2" }`)

	tests := []struct {
		name         string
		body         string
		expectedCode int
	}{
		{
			name:         "create muting",
			body:         `{ "target_user_id": "` + targetUserID + `" }`,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid JSON body",
			body:         `{ "target_user_id": "` + targetUserID,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid body",
			body:         `{ "invalid": "test" }`,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "duplicated muting",
			body:         `{ "target_user_id": "` + targetUserID + `" }`,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest(
			"POST",
			"/api/users/{id}/muting",
			strings.NewReader(test.body),
		)
		req.SetPathValue("id", sourceUserID)

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
	userID := s.newTestUser(`{ "username": "test", "display_name": "test" }`)
	postID := s.newTestPost(fmt.Sprintf(`{ "user_id": "%s", "text": "test" }`, userID))

	tests := []struct {
		name         string
		body         string
		expectedCode int
	}{
		{
			name:         "create repost",
			body:         fmt.Sprintf(`{ "post_id": "%s", "user_id": "%s" }`, postID, userID),
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid JSON body",
			body:         fmt.Sprintf(`{ "post_id": "%s", "user_id": "%s }`, postID, userID),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid body",
			body:         fmt.Sprintf(`{ "post_id": "%s" }`, postID),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "non-existent user id",
			body:         fmt.Sprintf(`{ "post_id": "%s", "user_id": "%s" }`, postID, uuid.New()),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "non-existent post id",
			body:         fmt.Sprintf(`{ "post_id": "%s", "user_id": "%s" }`, uuid.New(), userID),
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "duplicated repost",
			body:         fmt.Sprintf(`{ "post_id": "%s", "user_id": "%s" }`, postID, userID),
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

func (s *HandlersTestSuite) TestDeleteRepost() {
	userID := s.newTestUser(`{ "username": "test", "display_name": "test" }`)
	postID := s.newTestPost(fmt.Sprintf(`{ "user_id": "%s", "text": "test" }`, userID))
	s.newTestRepost(userID, postID)

	tests := []struct {
		name         string
		expectedCode int
	}{
		{
			name:         "delete repost",
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "non-existent repost",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(
			"DELETE",
			"/api/posts/reposts/{user_id}/{post_id}",
			strings.NewReader(""),
		)
		req.SetPathValue("user_id", userID)
		req.SetPathValue("post_id", postID)

		DeleteRepost(rr, req, s.db)

		if rr.Code != test.expectedCode {
			s.T().Errorf("%s: wrong code returned; expected %d, but got %d", test.name, test.expectedCode, rr.Code)
		}

	}

}

func (s *HandlersTestSuite) TestGetUserPostsTimeline() {
	// This test method verifies the number of posts in the response body.
	user1ID := s.newTestUser(`{ "username": "test1", "display_name": "test1" }`)
	_ = s.newTestPost(fmt.Sprintf(`{ "user_id": "%s", "text": "test1" }`, user1ID))
	user2ID := s.newTestUser(`{ "username": "test2", "display_name": "test2" }`)

	tests := []struct {
		name          string
		userID        string
		expectedCount int
	}{
		{
			name:          "get user posts",
			userID:        user1ID,
			expectedCount: 1,
		},
		{
			name:          "get no posts",
			userID:        user2ID,
			expectedCount: 0,
		},
	}

	for _, test := range tests {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(
			"GET",
			"/api/users/{id}/posts",
			strings.NewReader(""),
		)
		req.SetPathValue("id", test.userID)

		GetUserPostsTimeline(rr, req, s.getSpecificUserPostsUsecase)
		var posts []*entities.Post

		decoder := json.NewDecoder(rr.Body)
		err := decoder.Decode(&posts)
		if err != nil {
			s.T().Errorf("%s: failed to decode response", test.name)
		}

		if len(posts) != test.expectedCount {
			s.T().Errorf("%s: wrong number of posts returned; expected %d, but got %d", test.name, test.expectedCount, len(posts))
		}
	}
}

func (s *HandlersTestSuite) TestGetReverseChronologicalHomeTimeline() {
	// This test method verifies the number of posts in the response body.
	user1ID := s.newTestUser(`{ "username": "test1", "display_name": "test1" }`)
	_ = s.newTestPost(fmt.Sprintf(`{ "user_id": "%s", "text": "test1" }`, user1ID))
	user2ID := s.newTestUser(`{ "username": "test2", "display_name": "test2" }`)
	_ = s.newTestPost(fmt.Sprintf(`{ "user_id": "%s", "text": "test2" }`, user2ID))
	user3ID := s.newTestUser(`{ "username": "test3", "display_name": "test3" }`)
	_ = s.newTestPost(fmt.Sprintf(`{ "user_id": "%s", "text": "test3" }`, user3ID))
	user4ID := s.newTestUser(`{ "username": "test4", "display_name": "test4" }`)
	s.newTestFollow(user3ID, user2ID)

	tests := []struct {
		name          string
		userID        string
		expectedCount int
	}{
		{
			name:          "get only a target user posts",
			userID:        user1ID,
			expectedCount: 1,
		},
		{
			name:          "get a target user and following users posts",
			userID:        user3ID,
			expectedCount: 2,
		},
		{
			name:          "get no posts",
			userID:        user4ID,
			expectedCount: 0,
		},
	}

	for _, test := range tests {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(
			"GET",
			"/api/users/{id}/timelines/reverse_chronological",
			strings.NewReader(""),
		)
		req.SetPathValue("id", test.userID)

		GetReverseChronologicalHomeTimeline(rr, req, s.getUserAndFolloweePostsUsecase)
		var posts []*entities.Post

		decoder := json.NewDecoder(rr.Body)
		err := decoder.Decode(&posts)
		if err != nil {
			s.T().Errorf("%s: failed to decode response", test.name)
		}

		if len(posts) != test.expectedCount {
			s.T().Errorf("%s: wrong number of posts returned; expected %d, but got %d", test.name, test.expectedCount, len(posts))
		}
	}
}

func (s *HandlersTestSuite) newTestRepost(userID, postID string) {
	req := httptest.NewRequest(
		"POST",
		"/api/posts/reposts",
		strings.NewReader(fmt.Sprintf(`{ "post_id": "%s", "user_id": "%s" }`, postID, userID)),
	)
	rr := httptest.NewRecorder()
	CreateRepost(rr, req, s.db)
}

func (s *HandlersTestSuite) newTestUser(body string) string {
	req := httptest.NewRequest(
		"POST",
		"/api/users",
		strings.NewReader(body),
	)
	rr := httptest.NewRecorder()
	CreateUser(rr, req, s.db)

	var user entities.User
	_ = json.NewDecoder(rr.Body).Decode(&user)
	sourceUserID := user.ID.String()
	return sourceUserID
}

func (s *HandlersTestSuite) newTestPost(body string) string {
	req := httptest.NewRequest(
		"POST",
		"/api/posts",
		strings.NewReader(body),
	)
	rr := httptest.NewRecorder()
	CreatePost(rr, req, s.db)

	var post entities.Post
	_ = json.NewDecoder(rr.Body).Decode(&post)
	postID := post.ID.String()
	return postID
}

func (s *HandlersTestSuite) newTestLike(userID string, postID string) {
	req := httptest.NewRequest(
		"POST",
		"/api/users/{id}/likes",
		strings.NewReader(fmt.Sprintf(`{ "post_id": "%s" }`, postID)),
	)
	req.SetPathValue("id", userID)

	rr := httptest.NewRecorder()
	LikePost(rr, req, s.db)
}

func (s *HandlersTestSuite) newTestFollow(sourceUserID string, targetUserID string) {
	req := httptest.NewRequest(
		"POST",
		"/api/users/{id}/following",
		strings.NewReader(fmt.Sprintf(`{ "target_user_id": "%s" }`, targetUserID)),
	)
	req.SetPathValue("id", sourceUserID)

	rr := httptest.NewRecorder()
	CreateFollowship(rr, req, s.db)
}

// TestHandlersTestSuite runs all of the tests attached to HandlersTestSuite.
func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}
