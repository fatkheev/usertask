package service

import (
	"testing"
	"usertask/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByID(userID int) (*models.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUserPoints(userID int, points int) error {
	args := m.Called(userID, points)
	return args.Error(0)
}

func (m *MockUserRepository) SetUserReferrer(userID int, referrerID int) error {
	args := m.Called(userID, referrerID)
	return args.Error(0)
}

func (m *MockUserRepository) CompleteTask(userID int, taskType string, points int) error {
	args := m.Called(userID, taskType, points)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserReferrer(userID int) (int, error) {
	args := m.Called(userID)
	return args.Int(0), args.Error(1)
}

func TestGetUserStatus(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &models.User{ID: 1, Username: "testuser", Points: 100}

	mockRepo.On("GetUserByID", 1).Return(user, nil)

	result, err := service.GetUserStatus(1)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testuser", result.Username)
}

func TestCompleteTask(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &models.User{ID: 1, Username: "testuser", Points: 100}

	mockRepo.On("GetUserByID", 1).Return(user, nil)
	mockRepo.On("CompleteTask", 1, "test_task", 50).Return(nil)

	mockRepo.On("GetUserReferrer", 1).Return(0, nil)

	err := service.CompleteTask(1, "test_task", 50)
	assert.NoError(t, err)

	mockRepo.AssertCalled(t, "CompleteTask", 1, "test_task", 50)
	mockRepo.AssertCalled(t, "GetUserReferrer", 1)
}

func TestCompleteTask_WithReferrerBonus(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &models.User{ID: 2, Username: "newuser", Points: 0}

	mockRepo.On("GetUserByID", 2).Return(user, nil)
	mockRepo.On("CompleteTask", 2, "test_task", 100).Return(nil)
	mockRepo.On("GetUserReferrer", 2).Return(1, nil)
	mockRepo.On("UpdateUserPoints", 1, 50).Return(nil)

	err := service.CompleteTask(2, "test_task", 100)
	assert.NoError(t, err)

	mockRepo.AssertCalled(t, "UpdateUserPoints", 1, 50)
}

func TestSetReferrer(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &models.User{ID: 2, Username: "user2", Points: 0, ReferrerID: nil}
	referrer := &models.User{ID: 1, Username: "user1", Points: 0}

	mockRepo.On("GetUserByID", 2).Return(user, nil)
	mockRepo.On("GetUserByID", 1).Return(referrer, nil)
	mockRepo.On("SetUserReferrer", 2, 1).Return(nil)
	mockRepo.On("UpdateUserPoints", 1, 50).Return(nil) 
	mockRepo.On("CompleteTask", 1, "referral", 50).Return(nil)

	err := service.SetReferrer(2, 1)
	assert.NoError(t, err)

	mockRepo.AssertCalled(t, "SetUserReferrer", 2, 1)
	mockRepo.AssertCalled(t, "UpdateUserPoints", 1, 50)
	mockRepo.AssertCalled(t, "CompleteTask", 1, "referral", 50)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	user := &models.User{ID: 3, Username: "newuser", Points: 0}

	mockRepo.On("CreateUser", "newuser").Return(user, nil)

	createdUser, token, err := service.CreateUser("newuser")
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, "newuser", createdUser.Username)
	assert.NotEmpty(t, token)
}

func (m *MockUserRepository) GetLeaderboard(limit int) ([]models.User, error) {
	args := m.Called(limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

func TestGetLeaderboard(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	users := []models.User{
		{ID: 1, Username: "Alice", Points: 200},
		{ID: 2, Username: "Bob", Points: 150},
	}

	mockRepo.On("GetLeaderboard", 2).Return(users, nil)

	result, err := service.GetLeaderboard(2)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Alice", result[0].Username)
}