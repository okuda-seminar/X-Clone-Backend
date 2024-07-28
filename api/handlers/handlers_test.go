package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"x-clone-backend/entities"

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

// TestHandlersTestSuite runs all of the tests attached to HandlersTestSuite.
func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}
