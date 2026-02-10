// Package domain はビジネスロジックの中核となるドメイン層を定義します
package domain

import (
	"crypto/rand"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

// User はユーザーエンティティを表します
type User struct {
	ID    string `json:"id" gorm:"primaryKey;type:varchar(26)"`
	Name  string `json:"name" gorm:"type:varchar(100)"`
	Email string `json:"email" gorm:"type:varchar(255);uniqueIndex"`
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
	if !IsValidEmail(trimmedEmail) {
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
