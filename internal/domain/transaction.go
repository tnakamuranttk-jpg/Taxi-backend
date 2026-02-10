// Package domain はビジネスロジックの中核となるドメイン層を定義します
package domain

// TransactionManager はトランザクション管理のインターフェースです
type TransactionManager interface {
	// ExecuteInTransaction はトランザクション内で関数を実行します
	// 関数がエラーを返した場合はロールバック、nilの場合はコミットされます
	ExecuteInTransaction(fn func() error) error
}
