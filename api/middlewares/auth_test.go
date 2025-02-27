package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"x-clone-backend/internal/app/services"

	"github.com/google/uuid"
)

// TestJWTMiddleware_TokenValidation verifies the basic functionality of JWTMiddleware.
// It checks if a valid JWT token grants access, an invalid token denies access,
// and requests missing the Authorization header are rejected.
func TestJWTMiddleware_TokenValidation(t *testing.T) {
	secretKey := "test_secret_key"
	authService := services.NewAuthService(secretKey)

	// Generate a valid JWT token
	tokenString, _ := authService.GenerateJWT(uuid.New(), "test_user")

	// Create a test handler to verify JWT claims in the request context
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(UserContextKey).(*services.UserClaims)
		if !ok {
			t.Error("Failed to retrieve claims from context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims.Username != "test_user" {
			t.Errorf("Expected username 'test_user', got '%v'", claims.Username)
		}
	})

	// Apply the JWT middleware
	handlerToTest := JWTMiddleware(authService)(testHandler)

	tests := map[string]struct {
		token          string
		expectedStatus int
	}{
		"Valid JWT Token": {
			token:          "Bearer " + tokenString,
			expectedStatus: http.StatusOK,
		},
		"Invalid JWT Token": {
			token:          "Bearer invalidtoken",
			expectedStatus: http.StatusUnauthorized,
		},
		"Missing Authorization Header": {
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			rr := httptest.NewRecorder()

			handlerToTest.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code for %v: got %v want %v", name, status, tt.expectedStatus)
			}
		})
	}
}

// TestJWTMiddleware_Endpoints verifies the behavior of JWTMiddleware for different API endpoints.
// It tests whether authentication is enforced correctly for protected routes
// and ensures public routes can be accessed without a JWT token.
func TestJWTMiddleware_Endpoints(t *testing.T) {
	secretKey := "test_secret_key"
	authService := services.NewAuthService(secretKey)

	// Generate a valid JWT token
	tokenString, _ := authService.GenerateJWT(uuid.New(), "test_user")

	// Create a new multiplexer with different API endpoints
	mux := http.NewServeMux()

	// Public endpoint that does not require authentication
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Protected endpoint that requires authentication
	mux.HandleFunc("/api/protected", func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(UserContextKey).(*services.UserClaims)
		if !ok {
			t.Error("Failed to retrieve claims from context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims.Username != "test_user" {
			t.Errorf("Expected username 'test_user', got '%v'", claims.Username)
		}

		w.WriteHeader(http.StatusOK)
	})

	// Apply the JWT middleware
	handlerToTest := JWTMiddleware(authService)(mux)

	tests := map[string]struct {
		method         string
		url            string
		token          string
		expectedStatus int
	}{
		"Public Endpoint - Without JWT": {
			method:         "POST",
			url:            "/api/users",
			token:          "",
			expectedStatus: http.StatusOK,
		},
		"Protected Endpoint - Valid JWT": {
			method:         "GET",
			url:            "/api/protected",
			token:          "Bearer " + tokenString,
			expectedStatus: http.StatusOK,
		},
		"Protected Endpoint - Without JWT": {
			method:         "GET",
			url:            "/api/protected",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
		"Protected Endpoint - Invalid JWT": {
			method:         "GET",
			url:            "/api/protected",
			token:          "Bearer invalid_token",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			rr := httptest.NewRecorder()

			handlerToTest.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Handler returned wrong status code for %v: got %v want %v", name, status, tt.expectedStatus)
			}
		})
	}
}
