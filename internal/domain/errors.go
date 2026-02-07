// Package domain はビジネスロジックの中核となるドメイン層を定義します
package domain

import "errors"

// アプリケーション全体で使用するカスタムエラー
var (
	// ErrUserNotFound はユーザーが見つからない場合のエラー
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists は既にユーザーが存在する場合のエラー
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrInvalidInput は入力値が不正な場合のエラー
	ErrInvalidInput = errors.New("invalid input")

	// ErrInvalidEmail はメールアドレスの形式が不正な場合のエラー
	ErrInvalidEmail = errors.New("invalid email format")

	// ErrEmptyName は名前が空の場合のエラー
	ErrEmptyName = errors.New("name cannot be empty")
)
