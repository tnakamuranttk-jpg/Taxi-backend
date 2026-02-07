// Package usecase はアプリケーションのビジネスロジックを実装します
package usecase

import (
	"github.com/tnakamura/taxi-backend/internal/domain"
)

// UserUsecase はユーザーに関するビジネスロジックを定義するインターフェースです
type UserUsecase interface {
	// GetUser は指定されたIDのユーザーを取得します
	GetUser(id string) (*domain.User, error)
	// CreateUser は新しいユーザーを作成します
	CreateUser(name, email string) (*domain.User, error)
	// GetAllUsers は全てのユーザーを取得します
	GetAllUsers() ([]*domain.User, error)
}

// userUsecase は UserUsecase インターフェースの具体的な実装です
type userUsecase struct {
	userRepo domain.UserRepository
}

// NewUserUsecase は新しい UserUsecase を生成します
func NewUserUsecase(repo domain.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: repo,
	}
}

// GetUser は指定されたIDのユーザーを取得します
func (u *userUsecase) GetUser(id string) (*domain.User, error) {
	if id == "" {
		return nil, domain.ErrInvalidInput
	}
	return u.userRepo.GetByID(id)
}

// CreateUser は新しいユーザーを作成して保存します
func (u *userUsecase) CreateUser(name, email string) (*domain.User, error) {
	// メールアドレスの重複チェック
	exists, err := u.userRepo.ExistsByEmail(email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrUserAlreadyExists
	}

	// ドメイン層の NewUser を呼び出して ULID を生成したユーザーオブジェクトを作成
	user, err := domain.NewUser(name, email)
	if err != nil {
		return nil, err
	}

	// リポジトリ経由でデータベースに保存
	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetAllUsers は登録されている全ユーザーを取得します
func (u *userUsecase) GetAllUsers() ([]*domain.User, error) {
	return u.userRepo.GetAll()
}
