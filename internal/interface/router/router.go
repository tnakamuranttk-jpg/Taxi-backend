// Package router はAPIのルーティング設定を管理します
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tnakamura/taxi-backend/internal/interface/handler"
)

// NewRouter は新しいGinルーターを作成し、ルーティングを設定します
func NewRouter(userHandler *handler.UserHandler, driverHandler *handler.DriverHandler) *gin.Engine {
	r := gin.Default()

	// ヘルスチェック
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 グループ
	v1 := r.Group("/api/v1")
	{
		// ユーザー関連のルーティング
		users := v1.Group("/users")
		{
			users.GET("", userHandler.GetAllUsers)
			users.GET("/:id", userHandler.GetUser)
			users.POST("", userHandler.CreateUser)
		}

		// ドライバー関連のルーティング
		drivers := v1.Group("/drivers")
		{
			drivers.GET("", driverHandler.GetAllDrivers)
			drivers.GET("/available", driverHandler.GetAvailableDrivers)
			drivers.GET("/:id", driverHandler.GetDriver)
			drivers.POST("", driverHandler.CreateDriver)
			drivers.PUT("/:id/location", driverHandler.UpdateDriverLocation)
			drivers.PUT("/:id/status", driverHandler.UpdateDriverStatus)
		}
	}

	return r
}
