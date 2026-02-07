// Package persistence はデータベースとの永続化処理を実装します
package persistence

import (
	"errors"

	"github.com/tnakamura/taxi-backend/internal/domain"
	"gorm.io/gorm"
)

// userRepository は domain.UserRepository の具体的な実装です
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository は新しい UserRepository を生成します
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

// GetByID は指定されたIDのユーザーを取得します
func (r *userRepository) GetByID(id string) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Create は新しいユーザーをデータベースに保存します
func (r *userRepository) Create(user *domain.User) error {
	if err := r.db.Create(user).Error; err != nil {
		// 重複エラーのチェック（MySQLのエラーコード1062）
		if isDuplicateError(err) {
			return domain.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

// GetAll は全てのユーザーを取得します
func (r *userRepository) GetAll() ([]*domain.User, error) {
	var users []*domain.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// ExistsByEmail は指定されたメールアドレスのユーザーが存在するかチェックします
func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	if err := r.db.Model(&domain.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// isDuplicateError はデータベースの重複エラーかどうかを判定します
func isDuplicateError(err error) bool {
	// エラーメッセージに "Duplicate entry" が含まれているかチェック
	return err != nil && (errors.Is(err, gorm.ErrDuplicatedKey) ||
		containsString(err.Error(), "Duplicate entry") ||
		containsString(err.Error(), "duplicate key"))
}

// containsString は文字列が部分文字列を含むかチェックします
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
