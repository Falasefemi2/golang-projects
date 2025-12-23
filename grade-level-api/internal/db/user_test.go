package db

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock DBExecutor
type MockDB struct {
	mock.Mock
}

func (m *MockDB) Exec(query string, args ...any) (sql.Result, error) {
	ret := m.Called(query, args)
	return ret.Get(0).(sql.Result), ret.Error(1)
}

// Mock sql.Result
type ResultMock struct {
	mock.Mock
}

func (r *ResultMock) LastInsertId() (int64, error) {
	args := r.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (r *ResultMock) RowsAffected() (int64, error) {
	args := r.Called()
	return args.Get(0).(int64), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	mockDB := new(MockDB)
	result := new(ResultMock)

	result.On("LastInertId").Return(int64(1), nil)
	mockDB.On("Exec",
		"INSERT INTO user (name, email, password, role) VALUES (?, ?, ?, ?)",
		mock.Anything, mock.Anything, mock.Anything, mock.Anything,
	).Return(result, nil)

	user, err := CreateUser(mockDB, "Femi", "femi@example.com", "secret123", Student)

	assert.NoError(t, err)
	assert.Equal(t, "Femi", user.Name)
	assert.Equal(t, "femi@example.com", user.Email)
	assert.Equal(t, string(Student), user.Role)
	assert.NotEmpty(t, user.Password)

	mockDB.AssertExpectations(t)
	result.AssertExpectations(t)
}
