package usecase_test

import (
	"errors"
	"testing"

	"github.com/UthitSawatdee/GoMarketAPI/internal/domain"
	usecase "github.com/UthitSawatdee/GoMarketAPI/internal/usecases"
)

// ==============================================
// MOCK IMPLEMENTATIONS
// ==============================================

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	users       map[string]*domain.User
	createError error
	getError    error
	updateError error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (m *MockUserRepository) Create(user *domain.User) error {
	if m.createError != nil {
		return m.createError
	}
	user.ID = uint(len(m.users) + 1)
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	if user, exists := m.users[email]; exists {
		return user, nil
	}
	return nil, nil
}

func (m *MockUserRepository) GetUserByID(id uint) (*domain.User, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) Update(user *domain.User) error {
	if m.updateError != nil {
		return m.updateError
	}
	if existing, exists := m.users[user.Email]; exists {
		existing.Username = user.Username
		existing.Password = user.Password
		return nil
	}
	return errors.New("user not found")
}

func (m *MockUserRepository) AllUsers() ([]*domain.User, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	var users []*domain.User
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

// MockPasswordService is a mock implementation of PasswordService
type MockPasswordService struct {
	hashError   error
	verifyError bool
}

func NewMockPasswordService() *MockPasswordService {
	return &MockPasswordService{}
}

func (m *MockPasswordService) Hash(password string) (string, error) {
	if m.hashError != nil {
		return "", m.hashError
	}
	return "hashed_" + password, nil
}

func (m *MockPasswordService) Verify(password, hashedPassword string) bool {
	if m.verifyError {
		return false
	}
	return "hashed_"+password == hashedPassword
}

// ==============================================
// USER SERVICE TESTS
// ==============================================

func TestUserService_CreateUser_Success(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()
	mockHash := NewMockPasswordService()
	service := usecase.NewUserService(mockRepo, mockHash)

	user := &domain.User{
		Email:    "test@example.com",
		Password: "password123",
		Username: "testuser",
		Role:     "user",
	}

	// Act
	err := service.CreateUser(user)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if user.ID == 0 {
		t.Error("Expected user ID to be set")
	}
	if user.Password != "hashed_password123" {
		t.Errorf("Expected password to be hashed, got: %s", user.Password)
	}
}

func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()
	mockHash := NewMockPasswordService()
	service := usecase.NewUserService(mockRepo, mockHash)

	// Create first user
	existingUser := &domain.User{
		Email:    "test@example.com",
		Password: "password123",
		Username: "existinguser",
	}
	mockRepo.users[existingUser.Email] = existingUser

	// Try to create duplicate
	newUser := &domain.User{
		Email:    "test@example.com",
		Password: "password456",
		Username: "newuser",
	}

	// Act
	err := service.CreateUser(newUser)

	// Assert
	if err == nil {
		t.Error("Expected error for duplicate email, got nil")
	}
	if err.Error() != "email already registered" {
		t.Errorf("Expected 'email already registered' error, got: %v", err)
	}
}

func TestUserService_CreateUser_HashError(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()
	mockHash := NewMockPasswordService()
	mockHash.hashError = errors.New("hash failed")
	service := usecase.NewUserService(mockRepo, mockHash)

	user := &domain.User{
		Email:    "test@example.com",
		Password: "password123",
		Username: "testuser",
	}

	// Act
	err := service.CreateUser(user)

	// Assert
	if err == nil {
		t.Error("Expected error from hash failure, got nil")
	}
}

func TestUserService_GetUserByID_Success(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()
	mockHash := NewMockPasswordService()
	service := usecase.NewUserService(mockRepo, mockHash)

	existingUser := &domain.User{
		ID:       1,
		Email:    "test@example.com",
		Username: "testuser",
	}
	mockRepo.users[existingUser.Email] = existingUser

	// Act
	user, err := service.GetUserByID(1)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if user == nil {
		t.Fatal("Expected user, got nil")
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got: %s", user.Email)
	}
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()
	mockHash := NewMockPasswordService()
	service := usecase.NewUserService(mockRepo, mockHash)

	// Act
	user, err := service.GetUserByID(999)

	// Assert
	if err == nil {
		t.Error("Expected error for non-existent user, got nil")
	}
	if user != nil {
		t.Error("Expected nil user for non-existent ID")
	}
}

func TestUserService_AllUsers_Success(t *testing.T) {
	// Arrange
	mockRepo := NewMockUserRepository()
	mockHash := NewMockPasswordService()
	service := usecase.NewUserService(mockRepo, mockHash)

	mockRepo.users["user1@example.com"] = &domain.User{ID: 1, Email: "user1@example.com"}
	mockRepo.users["user2@example.com"] = &domain.User{ID: 2, Email: "user2@example.com"}

	// Act
	users, err := service.AllUsers()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got: %d", len(users))
	}
}
