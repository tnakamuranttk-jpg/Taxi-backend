// Package dto はAPIのリクエスト/レスポンス用のデータ構造を定義します
package dto

import (
	"time"

	"github.com/tnakamura/taxi-backend/internal/domain"
)

// CreateDriverRequest はドライバー作成リクエストの構造体です
type CreateDriverRequest struct {
	Name          string `json:"name" binding:"required,min=1,max=100"`
	Email         string `json:"email" binding:"required,email"`
	LicenseNumber string `json:"license_number" binding:"required,min=1,max=50"`
}

// UpdateDriverLocationRequest はドライバー位置更新リクエストの構造体です
type UpdateDriverLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// UpdateDriverStatusRequest はドライバーステータス更新リクエストの構造体です
type UpdateDriverStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=available busy offline"`
}

// DriverResponse はドライバー情報のレスポンス構造体です
type DriverResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	LicenseNumber string    `json:"license_number"`
	Status        string    `json:"status"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ToDriverResponse はドメインモデルからレスポンス用DTOに変換します
func ToDriverResponse(driver *domain.Driver) *DriverResponse {
	if driver == nil {
		return nil
	}
	return &DriverResponse{
		ID:            driver.ID,
		Name:          driver.Name,
		Email:         driver.Email,
		LicenseNumber: driver.LicenseNumber,
		Status:        string(driver.Status),
		Latitude:      driver.Latitude,
		Longitude:     driver.Longitude,
		CreatedAt:     driver.CreatedAt,
		UpdatedAt:     driver.UpdatedAt,
	}
}

// ToDriverResponseList は複数のドメインモデルからレスポンス用DTOのリストに変換します
func ToDriverResponseList(drivers []*domain.Driver) []*DriverResponse {
	result := make([]*DriverResponse, len(drivers))
	for i, driver := range drivers {
		result[i] = ToDriverResponse(driver)
	}
	return result
}
