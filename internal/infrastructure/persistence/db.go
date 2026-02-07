// Package persistence はデータベースとの永続化処理を実装します
package persistence

import (
	"github.com/tnakamura/taxi-backend/internal/config"
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
