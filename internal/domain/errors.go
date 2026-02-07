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

	// ErrDriverNotFound はドライバーが見つからない場合のエラー
	ErrDriverNotFound = errors.New("driver not found")

	// ErrDriverAlreadyExists は既にドライバーが存在する場合のエラー
	ErrDriverAlreadyExists = errors.New("driver already exists")

	// ErrEmptyLicenseNumber は免許番号が空の場合のエラー
	ErrEmptyLicenseNumber = errors.New("license number cannot be empty")

	// ErrInvalidLatitude は緯度が不正な場合のエラー
	ErrInvalidLatitude = errors.New("latitude must be between -90 and 90")

	// ErrInvalidLongitude は経度が不正な場合のエラー
	ErrInvalidLongitude = errors.New("longitude must be between -180 and 180")

	// ErrInvalidDriverStatus はドライバーステータスが不正な場合のエラー
	ErrInvalidDriverStatus = errors.New("invalid driver status")
)
