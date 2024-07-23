package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
			name:         "invalid body",
			body:         `{ "invalid": "test" }`,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "duplicated username",
			body:         `{ "username": "test", "display_name": "duplicated" }`,
			expectedCode: http.StatusInternalServerError,
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

// TestHandlersTestSuite runs all of the tests attached to HandlersTestSuite.
func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}
