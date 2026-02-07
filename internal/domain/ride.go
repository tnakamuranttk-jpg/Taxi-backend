// Package domain はビジネスロジックの中核となるドメイン層を定義します
package domain

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// RideStatus は配車のステータスを表します
type RideStatus string

const (
	// RideStatusRequested は配車依頼中
	RideStatusRequested RideStatus = "requested"
	// RideStatusAccepted はドライバーが承諾
	RideStatusAccepted RideStatus = "accepted"
	// RideStatusArrived はドライバーが乗車地点に到着
	RideStatusArrived RideStatus = "arrived"
	// RideStatusOngoing は乗車中（走行中）
	RideStatusOngoing RideStatus = "ongoing"
	// RideStatusCompleted は乗車完了
	RideStatusCompleted RideStatus = "completed"
	// RideStatusCancelled はキャンセル
	RideStatusCancelled RideStatus = "cancelled"
)

// Ride は配車リクエストエンティティを表します
type Ride struct {
	ID               string     `json:"id" gorm:"primaryKey;type:varchar(26)"`
	PassengerID      string     `json:"passenger_id" gorm:"type:varchar(26);index"`
	DriverID         string     `json:"driver_id" gorm:"type:varchar(26);index"`
	PickupLatitude   float64    `json:"pickup_latitude"`
	PickupLongitude  float64    `json:"pickup_longitude"`
	DropoffLatitude  float64    `json:"dropoff_latitude"`
	DropoffLongitude float64    `json:"dropoff_longitude"`
	Status           RideStatus `json:"status" gorm:"type:varchar(20);default:requested;index"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// RideRepository は配車データの永続化に関するインターフェースです
type RideRepository interface {
	GetByID(id string) (*Ride, error)
	Create(ride *Ride) error
	Update(ride *Ride) error
	GetByPassengerID(passengerID string) ([]*Ride, error)
	GetByDriverID(driverID string) ([]*Ride, error)
	GetActiveByPassengerID(passengerID string) (*Ride, error)
	GetActiveByDriverID(driverID string) (*Ride, error)
}

// NewRide は新しい配車リクエストを作成します（バリデーション付き）
func NewRide(passengerID string, pickupLat, pickupLng, dropoffLat, dropoffLng float64) (*Ride, error) {
	// 乗客IDのバリデーション
	if passengerID == "" {
		return nil, ErrInvalidInput
	}

	// 緯度経度のバリデーション
	if pickupLat < -90 || pickupLat > 90 {
		return nil, ErrInvalidLatitude
	}
	if pickupLng < -180 || pickupLng > 180 {
		return nil, ErrInvalidLongitude
	}
	if dropoffLat < -90 || dropoffLat > 90 {
		return nil, ErrInvalidLatitude
	}
	if dropoffLng < -180 || dropoffLng > 180 {
		return nil, ErrInvalidLongitude
	}

	// ULID の生成
	entropy := ulid.Monotonic(rand.Reader, 0)
	id, err := ulid.New(ulid.Timestamp(time.Now()), entropy)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Ride{
		ID:               id.String(),
		PassengerID:      passengerID,
		DriverID:         "", // マッチング前は空
		PickupLatitude:   pickupLat,
		PickupLongitude:  pickupLng,
		DropoffLatitude:  dropoffLat,
		DropoffLongitude: dropoffLng,
		Status:           RideStatusRequested,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// Accept はドライバーが配車リクエストを承諾します
func (r *Ride) Accept(driverID string) error {
	if r.Status != RideStatusRequested {
		return ErrInvalidRideStatus
	}
	if driverID == "" {
		return ErrInvalidInput
	}

	r.DriverID = driverID
	r.Status = RideStatusAccepted
	r.UpdatedAt = time.Now()
	return nil
}

// Arrive はドライバーが乗車地点に到着したことを記録します
func (r *Ride) Arrive() error {
	if r.Status != RideStatusAccepted {
		return ErrInvalidRideStatus
	}

	r.Status = RideStatusArrived
	r.UpdatedAt = time.Now()
	return nil
}

// Start は乗車を開始します
func (r *Ride) Start() error {
	if r.Status != RideStatusArrived {
		return ErrInvalidRideStatus
	}

	r.Status = RideStatusOngoing
	r.UpdatedAt = time.Now()
	return nil
}

// Complete は乗車を完了します
func (r *Ride) Complete() error {
	if r.Status != RideStatusOngoing {
		return ErrInvalidRideStatus
	}

	r.Status = RideStatusCompleted
	r.UpdatedAt = time.Now()
	return nil
}

// Cancel は配車をキャンセルします
func (r *Ride) Cancel() error {
	// 完了済みはキャンセル不可
	if r.Status == RideStatusCompleted || r.Status == RideStatusCancelled {
		return ErrInvalidRideStatus
	}

	r.Status = RideStatusCancelled
	r.UpdatedAt = time.Now()
	return nil
}

// IsActive は配車が進行中かどうかを返します
func (r *Ride) IsActive() bool {
	switch r.Status {
	case RideStatusRequested, RideStatusAccepted, RideStatusArrived, RideStatusOngoing:
		return true
	default:
		return false
	}
}
