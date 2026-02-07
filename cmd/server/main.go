package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tnakamura/taxi-backend/internal/config"
	"github.com/tnakamura/taxi-backend/internal/domain"
	"github.com/tnakamura/taxi-backend/internal/infrastructure/persistence"
	"github.com/tnakamura/taxi-backend/internal/interface/handler"
	"github.com/tnakamura/taxi-backend/internal/interface/router"
	"github.com/tnakamura/taxi-backend/internal/usecase"
)

func main() {
	// 0. 環境変数のロード
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 1. 設定のロード
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. データベース接続の初期化
	db, err := persistence.NewDB(&cfg.DB)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	// 3. 自動マイグレーション (開発中のみ推奨)
	if err := db.AutoMigrate(&domain.User{}, &domain.Driver{}, &domain.Ride{}); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// 4. 依存関係の注入 (DI)
	// User
	userRepo := persistence.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	// Driver
	driverRepo := persistence.NewDriverRepository(db)
	driverUsecase := usecase.NewDriverUsecase(driverRepo)
	driverHandler := handler.NewDriverHandler(driverUsecase)

	// Ride
	rideRepo := persistence.NewRideRepository(db)
	rideUsecase := usecase.NewRideUsecase(rideRepo, driverRepo, userRepo)
	rideHandler := handler.NewRideHandler(rideUsecase)

	// 5. ルーターの設定
	r := router.NewRouter(userHandler, driverHandler, rideHandler)

	// 6. サーバー起動
	log.Printf("Server starting on :%s...", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
