// Package domain はビジネスロジックの中核となるドメイン層を定義します
package domain

import (
	"crypto/rand"
	"regexp"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

// User はユーザーエンティティを表します
type User struct {
	ID    string `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email" gorm:"uniqueIndex"`
}

// UserRepository はユーザーデータの永続化に関するインターフェースです
type UserRepository interface {
	GetByID(id string) (*User, error)
	Create(user *User) error
	GetAll() ([]*User, error)
	ExistsByEmail(email string) (bool, error)
}

// NewUser は新しいユーザーを作成します（バリデーション付き）
func NewUser(name, email string) (*User, error) {
	// 名前のバリデーション
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return nil, ErrEmptyName
	}

	// メールアドレスのバリデーション
	trimmedEmail := strings.TrimSpace(email)
	if !isValidEmail(trimmedEmail) {
		return nil, ErrInvalidEmail
	}

	// ULID の生成
	entropy := ulid.Monotonic(rand.Reader, 0)
	id, err := ulid.New(ulid.Timestamp(time.Now()), entropy)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:    id.String(),
		Name:  trimmedName,
		Email: trimmedEmail,
	}, nil
}

// isValidEmail はメールアドレスの形式をチェックします
func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	// 簡易的なメールアドレス形式チェック
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
