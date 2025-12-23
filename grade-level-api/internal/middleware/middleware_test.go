package middleware_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/falasefemi2/gradesystem/internal/db"
	"github.com/falasefemi2/gradesystem/internal/middleware"
	"github.com/falasefemi2/gradesystem/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

func generateTestJWT(email string, secret []byte) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &middleware.Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Exec(query string, args ...any) (sql.Result, error) {
	ret := m.Called(query, args)
	return ret.Get(0).(sql.Result), ret.Error(1)
}

func (m *MockDB) QueryRow(query string, args ...any) *sql.Row {
	return nil
}

func TestRoleAuth(t *testing.T) {
	// Create a mock handler that will be protected by the middleware
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a request with a valid token for an admin user
	adminToken, err := generateTestJWT("admin@example.com", []byte("supersecretkey"))
	if err != nil {
		t.Fatalf("Failed to generate admin token: %v", err)
	}

	// Create a request with a valid token for a non-admin user
	userToken, err := generateTestJWT("user@example.com", []byte("supersecretkey"))
	if err != nil {
		t.Fatalf("Failed to generate user token: %v", err)
	}

	// Create a mock for GetUserByEmail
	originalGetUserByEmail := db.GetUserByEmail
	defer func() { db.GetUserByEmail = originalGetUserByEmail }()

	db.GetUserByEmail = func(email string) (*models.User, error) {
		if email == "admin@example.com" {
			return &models.User{Email: email, Role: "admin"}, nil
		}
		return &models.User{Email: email, Role: "student"}, nil
	}

	testCases := []struct {
		name           string
		token          string
		requiredRole   db.Role
		expectedStatus int
	}{
		{
			name:           "Admin Access",
			token:          adminToken,
			requiredRole:   db.Admin,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-Admin Access",
			token:          userToken,
			requiredRole:   db.Admin,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "No Token",
			token:          "",
			requiredRole:   db.Admin,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/admin/users", nil)
			if err != nil {
				t.Fatal(err)
			}

			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}

			rr := httptest.NewRecorder()
			handler := middleware.RoleAuth(mockHandler, tc.requiredRole)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}
		})
	}
}
