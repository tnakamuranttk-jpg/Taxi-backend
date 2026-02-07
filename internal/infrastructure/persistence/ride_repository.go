// Package persistence はデータベースとの永続化処理を実装します
package persistence

import (
	"errors"

	"github.com/tnakamura/taxi-backend/internal/domain"
	"gorm.io/gorm"
)

// rideRepository は domain.RideRepository の具体的な実装です
type rideRepository struct {
	db *gorm.DB
}

// NewRideRepository は新しい RideRepository を生成します
func NewRideRepository(db *gorm.DB) domain.RideRepository {
	return &rideRepository{
		db: db,
	}
}

// GetByID は指定されたIDの配車を取得します
func (r *rideRepository) GetByID(id string) (*domain.Ride, error) {
	var ride domain.Ride
	if err := r.db.First(&ride, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRideNotFound
		}
		return nil, err
	}
	return &ride, nil
}

// Create は新しい配車をデータベースに保存します
func (r *rideRepository) Create(ride *domain.Ride) error {
	if err := r.db.Create(ride).Error; err != nil {
		return err
	}
	return nil
}

// Update は配車情報を更新します
func (r *rideRepository) Update(ride *domain.Ride) error {
	result := r.db.Save(ride)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrRideNotFound
	}
	return nil
}

// GetByPassengerID は指定された乗客IDの配車履歴を取得します
func (r *rideRepository) GetByPassengerID(passengerID string) ([]*domain.Ride, error) {
	var rides []*domain.Ride
	if err := r.db.Where("passenger_id = ?", passengerID).Order("created_at DESC").Find(&rides).Error; err != nil {
		return nil, err
	}
	return rides, nil
}

// GetByDriverID は指定されたドライバーIDの配車履歴を取得します
func (r *rideRepository) GetByDriverID(driverID string) ([]*domain.Ride, error) {
	var rides []*domain.Ride
	if err := r.db.Where("driver_id = ?", driverID).Order("created_at DESC").Find(&rides).Error; err != nil {
		return nil, err
	}
	return rides, nil
}

// GetActiveByPassengerID は指定された乗客の進行中の配車を取得します
func (r *rideRepository) GetActiveByPassengerID(passengerID string) (*domain.Ride, error) {
	var ride domain.Ride
	activeStatuses := []domain.RideStatus{
		domain.RideStatusRequested,
		domain.RideStatusAccepted,
		domain.RideStatusArrived,
		domain.RideStatusOngoing,
	}

	if err := r.db.Where("passenger_id = ? AND status IN ?", passengerID, activeStatuses).First(&ride).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRideNotFound
		}
		return nil, err
	}
	return &ride, nil
}

// GetActiveByDriverID は指定されたドライバーの進行中の配車を取得します
func (r *rideRepository) GetActiveByDriverID(driverID string) (*domain.Ride, error) {
	var ride domain.Ride
	activeStatuses := []domain.RideStatus{
		domain.RideStatusRequested,
		domain.RideStatusAccepted,
		domain.RideStatusArrived,
		domain.RideStatusOngoing,
	}

	if err := r.db.Where("driver_id = ? AND status IN ?", driverID, activeStatuses).First(&ride).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRideNotFound
		}
		return nil, err
	}
	return &ride, nil
}
