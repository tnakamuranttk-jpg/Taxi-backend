// Package dto はAPIのリクエスト/レスポンス用のデータ構造を定義します
package dto

import "github.com/tnakamura/taxi-backend/internal/domain"

// CreateUserRequest はユーザー作成リクエストの構造体です
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=100"`
	Email string `json:"email" binding:"required,email"`
}

// UserResponse はユーザー情報のレスポンス構造体です
type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ToUserResponse はドメインモデルからレスポンス用DTOに変換します
func ToUserResponse(user *domain.User) *UserResponse {
	if user == nil {
		return nil
	}
	return &UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

// ToUserResponseList は複数のドメインモデルからレスポンス用DTOのリストに変換します
func ToUserResponseList(users []*domain.User) []*UserResponse {
	result := make([]*UserResponse, len(users))
	for i, user := range users {
		result[i] = ToUserResponse(user)
	}
	return result
}

// ErrorResponse はエラーレスポンスの構造体です
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
