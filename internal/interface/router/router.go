// Package router はAPIのルーティング設定を管理します
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tnakamura/taxi-backend/internal/interface/handler"
)

// NewRouter は新しいGinルーターを作成し、ルーティングを設定します
func NewRouter(userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()

	// ヘルスチェック
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// ユーザー関連のルーティング
	r.GET("/users", userHandler.GetAllUsers)
	r.GET("/users/:id", userHandler.GetUser)
	r.POST("/users", userHandler.CreateUser)

	return r
}
