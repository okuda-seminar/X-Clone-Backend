package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

func (s *HandlersTestSuite) TestCreateUser() {
	tests := []struct {
		name         string
		body         string
		expectedCode int
	}{
		{
			name:         "create user",
			body:         `{ "username": "test", "display_name": "test", "password": "securepassword" }`,
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
			body:         `{ "username": "test", "display_name": "duplicated", "password": "securepassword" }`,
			expectedCode: http.StatusConflict,
		},
	}

	for _, test := range tests {
		createUserHandler := NewCreateUserHandler(s.db)

		req := httptest.NewRequest("POST", "/api/users", strings.NewReader(test.body))
		rr := httptest.NewRecorder()

		createUserHandler.CreateUser(rr, req)

		if rr.Code != test.expectedCode {
			s.T().Errorf("%s: wrong code returned; expected %d, but got %d", test.name, test.expectedCode, rr.Code)
		}
	}
}
