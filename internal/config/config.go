// Package config はアプリケーションの設定を管理します
package config

import (
	"fmt"
	"os"
)

// Config はアプリケーション全体の設定を保持します
type Config struct {
	DB     DBConfig
	Server ServerConfig
}

// DBConfig はデータベース接続に必要な設定を保持します
type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

// ServerConfig はサーバー設定を保持します
type ServerConfig struct {
	Port string
}

// Load は環境変数から設定を読み込みます
func Load() (*Config, error) {
	dbConfig := DBConfig{
		User:     getEnvOrDefault("DB_USER", "root"),
		Password: getEnvOrDefault("DB_PASSWORD", ""),
		Host:     getEnvOrDefault("DB_HOST", "127.0.0.1"),
		Port:     getEnvOrDefault("DB_PORT", "3306"),
		Name:     getEnvOrDefault("DB_NAME", "my_database"),
	}

	serverConfig := ServerConfig{
		Port: getEnvOrDefault("SERVER_PORT", "8080"),
	}

	// 必須項目のバリデーション
	if dbConfig.User == "" {
		return nil, fmt.Errorf("DB_USER is required")
	}
	if dbConfig.Name == "" {
		return nil, fmt.Errorf("DB_NAME is required")
	}

	return &Config{
		DB:     dbConfig,
		Server: serverConfig,
	}, nil
}

// DSN はデータベース接続文字列を生成します
func (c *DBConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name)
}

// getEnvOrDefault は環境変数を取得し、存在しない場合はデフォルト値を返します
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
