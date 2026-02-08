// Package domain はビジネスロジックの中核となるドメイン層を定義します
package domain

import "regexp"

// emailRegex はメールアドレスの形式チェック用の正規表現です
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// IsValidEmail はメールアドレスの形式をチェックします
func IsValidEmail(email string) bool {
	if email == "" {
		return false
	}
	return emailRegex.MatchString(email)
}
