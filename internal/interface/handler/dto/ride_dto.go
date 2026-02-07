// Package dto はAPIのリクエスト/レスポンス用のデータ構造を定義します
package dto

import (
	"time"

	"github.com/tnakamura/taxi-backend/internal/domain"
)

// CreateRideRequest は配車リクエスト作成の構造体です
type CreateRideRequest struct {
	PassengerID      string  `json:"passenger_id" binding:"required"`
	PickupLatitude   float64 `json:"pickup_latitude" binding:"required"`
	PickupLongitude  float64 `json:"pickup_longitude" binding:"required"`
	DropoffLatitude  float64 `json:"dropoff_latitude" binding:"required"`
	DropoffLongitude float64 `json:"dropoff_longitude" binding:"required"`
}

// AcceptRideRequest はドライバーが配車を承諾するリクエストの構造体です
type AcceptRideRequest struct {
	DriverID string `json:"driver_id" binding:"required"`
}

// RideResponse は配車情報のレスポンス構造体です
type RideResponse struct {
	ID               string    `json:"id"`
	PassengerID      string    `json:"passenger_id"`
	DriverID         string    `json:"driver_id,omitempty"`
	PickupLatitude   float64   `json:"pickup_latitude"`
	PickupLongitude  float64   `json:"pickup_longitude"`
	DropoffLatitude  float64   `json:"dropoff_latitude"`
	DropoffLongitude float64   `json:"dropoff_longitude"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// ToRideResponse はドメインモデルからレスポンス用DTOに変換します
func ToRideResponse(ride *domain.Ride) *RideResponse {
	if ride == nil {
		return nil
	}
	return &RideResponse{
		ID:               ride.ID,
		PassengerID:      ride.PassengerID,
		DriverID:         ride.DriverID,
		PickupLatitude:   ride.PickupLatitude,
		PickupLongitude:  ride.PickupLongitude,
		DropoffLatitude:  ride.DropoffLatitude,
		DropoffLongitude: ride.DropoffLongitude,
		Status:           string(ride.Status),
		CreatedAt:        ride.CreatedAt,
		UpdatedAt:        ride.UpdatedAt,
	}
}

// ToRideResponseList は複数のドメインモデルからレスポンス用DTOのリストに変換します
func ToRideResponseList(rides []*domain.Ride) []*RideResponse {
	result := make([]*RideResponse, len(rides))
	for i, ride := range rides {
		result[i] = ToRideResponse(ride)
	}
	return result
}
