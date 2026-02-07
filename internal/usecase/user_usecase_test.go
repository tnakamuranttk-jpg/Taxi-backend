package usecase_test

import (
	"testing"

	"github.com/tnakamura/taxi-backend/internal/domain"
	"github.com/tnakamura/taxi-backend/internal/usecase"
)

// MockUserRepository はテスト用のモックリポジトリです
type MockUserRepository struct {
	users       map[string]*domain.User
	shouldError bool
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (m *MockUserRepository) GetByID(id string) (*domain.User, error) {
	if m.shouldError {
		return nil, domain.ErrUserNotFound
	}
	user, ok := m.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) Create(user *domain.User) error {
	if m.shouldError {
		return domain.ErrUserAlreadyExists
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) GetAll() ([]*domain.User, error) {
	if m.shouldError {
		return nil, domain.ErrInvalidInput
	}
	result := make([]*domain.User, 0, len(m.users))
	for _, user := range m.users {
		result = append(result, user)
	}
	return result, nil
}

func (m *MockUserRepository) ExistsByEmail(email string) (bool, error) {
	for _, user := range m.users {
		if user.Email == email {
			return true, nil
		}
	}
	return false, nil
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	uc := usecase.NewUserUsecase(mockRepo)

	user, err := uc.CreateUser("テスト太郎", "test@example.com")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("expected user to be created")
	}

	if user.ID == "" {
		t.Error("expected user ID to be generated (ULID)")
	}

	if user.Name != "テスト太郎" {
		t.Errorf("expected name 'テスト太郎', got '%s'", user.Name)
	}

	if user.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", user.Email)
	}
}

func TestCreateUser_EmptyName(t *testing.T) {
	mockRepo := NewMockUserRepository()
	uc := usecase.NewUserUsecase(mockRepo)

	_, err := uc.CreateUser("", "test@example.com")

	if err != domain.ErrEmptyName {
		t.Errorf("expected ErrEmptyName, got %v", err)
	}
}

func TestCreateUser_InvalidEmail(t *testing.T) {
	mockRepo := NewMockUserRepository()
	uc := usecase.NewUserUsecase(mockRepo)

	_, err := uc.CreateUser("テスト太郎", "invalid-email")

	if err != domain.ErrInvalidEmail {
		t.Errorf("expected ErrInvalidEmail, got %v", err)
	}
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	mockRepo := NewMockUserRepository()
	uc := usecase.NewUserUsecase(mockRepo)

	// 最初のユーザーを作成
	_, err := uc.CreateUser("テスト太郎", "test@example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 同じメールアドレスで2人目を作成しようとする
	_, err = uc.CreateUser("テスト次郎", "test@example.com")
	if err != domain.ErrUserAlreadyExists {
		t.Errorf("expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestGetUser_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	uc := usecase.NewUserUsecase(mockRepo)

	// ユーザーを作成
	createdUser, _ := uc.CreateUser("テスト太郎", "test@example.com")

	// 取得
	user, err := uc.GetUser(createdUser.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID != createdUser.ID {
		t.Errorf("expected ID '%s', got '%s'", createdUser.ID, user.ID)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	mockRepo := NewMockUserRepository()
	uc := usecase.NewUserUsecase(mockRepo)

	_, err := uc.GetUser("nonexistent-id")

	if err != domain.ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestGetUser_EmptyID(t *testing.T) {
	mockRepo := NewMockUserRepository()
	uc := usecase.NewUserUsecase(mockRepo)

	_, err := uc.GetUser("")

	if err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestGetAllUsers_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	uc := usecase.NewUserUsecase(mockRepo)

	// ユーザーを複数作成
	uc.CreateUser("ユーザー1", "user1@example.com")
	uc.CreateUser("ユーザー2", "user2@example.com")

	users, err := uc.GetAllUsers()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestGetAllUsers_Empty(t *testing.T) {
	mockRepo := NewMockUserRepository()
	uc := usecase.NewUserUsecase(mockRepo)

	users, err := uc.GetAllUsers()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}
}
