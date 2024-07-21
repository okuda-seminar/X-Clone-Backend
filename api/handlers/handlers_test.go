package handlers

import (
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
)

// TODO: https://github.com/okuda-seminar/X-Clone-Backend/issues/63
// - [Users] Implement TestCreateUser
func (s *HandlersTestSuite) TestCreateUser() {
	log.Println("TestCreateUser was called.")
}

// TestHandlersTestSuite runs all of the tests attached to HandlersTestSuite.
func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}
