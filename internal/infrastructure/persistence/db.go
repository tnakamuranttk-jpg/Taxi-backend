// Package persistence はデータベースとの永続化処理を実装します
package persistence

import (
	"github.com/tnakamura/taxi-backend/internal/config"
	"github.com/tnakamura/taxi-backend/internal/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewDB はデータベース接続を初期化します
func NewDB(cfg *config.DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 接続確認 (Ping)
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// TransactionManager はトランザクション管理の実装です
type TransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager は新しい TransactionManager を生成します
func NewTransactionManager(db *gorm.DB) domain.TransactionManager {
	return &TransactionManager{db: db}
}

// ExecuteInTransaction はトランザクション内で関数を実行します
func (tm *TransactionManager) ExecuteInTransaction(fn func() error) error {
	return tm.db.Transaction(func(tx *gorm.DB) error {
		return fn()
	})
}
