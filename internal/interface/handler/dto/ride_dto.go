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

// EstimateRideRequest は料金見積もりリクエストの構造体です
type EstimateRideRequest struct {
	PickupLatitude   float64 `json:"pickup_latitude" binding:"required"`
	PickupLongitude  float64 `json:"pickup_longitude" binding:"required"`
	DropoffLatitude  float64 `json:"dropoff_latitude" binding:"required"`
	DropoffLongitude float64 `json:"dropoff_longitude" binding:"required"`
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
	FareAmount       int       `json:"fare_amount"`
	DistanceKm       float64   `json:"distance_km"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// FareEstimateResponse は料金見積もりのレスポンス構造体です
type FareEstimateResponse struct {
	BaseFare     int     `json:"base_fare"`
	DistanceFare int     `json:"distance_fare"`
	TotalFare    int     `json:"total_fare"`
	DistanceKm   float64 `json:"distance_km"`
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
		FareAmount:       ride.FareAmount,
		DistanceKm:       ride.DistanceKm,
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

// ToFareEstimateResponse はドメインの料金内訳をDTOに変換します
func ToFareEstimateResponse(fare *domain.FareBreakdown) *FareEstimateResponse {
	if fare == nil {
		return nil
	}
	return &FareEstimateResponse{
		BaseFare:     fare.BaseFare,
		DistanceFare: fare.DistanceFare,
		TotalFare:    fare.TotalFare,
		DistanceKm:   fare.DistanceKm,
	}
}
