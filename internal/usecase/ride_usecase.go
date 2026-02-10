// Package usecase はアプリケーションのビジネスロジックを実装します
package usecase

import (
	"github.com/tnakamura/taxi-backend/internal/domain"
)

// RideUsecase は配車に関するビジネスロジックを定義するインターフェースです
type RideUsecase interface {
	// CreateRide は新しい配車リクエストを作成します
	CreateRide(passengerID string, pickupLat, pickupLng, dropoffLat, dropoffLng float64) (*domain.Ride, error)
	// GetRide は指定されたIDの配車を取得します
	GetRide(id string) (*domain.Ride, error)
	// AcceptRide はドライバーが配車を承諾します
	AcceptRide(rideID, driverID string) (*domain.Ride, error)
	// ArriveAtPickup はドライバーが乗車地点に到着したことを記録します
	ArriveAtPickup(rideID string) (*domain.Ride, error)
	// StartRide は乗車を開始します
	StartRide(rideID string) (*domain.Ride, error)
	// CompleteRide は乗車を完了します
	CompleteRide(rideID string) (*domain.Ride, error)
	// CancelRide は配車をキャンセルします
	CancelRide(rideID string) (*domain.Ride, error)
	// GetPassengerRides は乗客の配車履歴を取得します
	GetPassengerRides(passengerID string) ([]*domain.Ride, error)
	// GetDriverRides はドライバーの配車履歴を取得します
	GetDriverRides(driverID string) ([]*domain.Ride, error)
	// EstimateFare は乗降車地点から料金を見積もります
	EstimateFare(pickupLat, pickupLng, dropoffLat, dropoffLng float64) (*domain.FareBreakdown, error)
}

// rideUsecase は RideUsecase インターフェースの具体的な実装です
type rideUsecase struct {
	rideRepo   domain.RideRepository
	driverRepo domain.DriverRepository
	userRepo   domain.UserRepository
	txManager  domain.TransactionManager
}

// NewRideUsecase は新しい RideUsecase を生成します
func NewRideUsecase(rideRepo domain.RideRepository, driverRepo domain.DriverRepository, userRepo domain.UserRepository, txManager domain.TransactionManager) RideUsecase {
	return &rideUsecase{
		rideRepo:   rideRepo,
		driverRepo: driverRepo,
		userRepo:   userRepo,
		txManager:  txManager,
	}
}

// CreateRide は新しい配車リクエストを作成します
func (u *rideUsecase) CreateRide(passengerID string, pickupLat, pickupLng, dropoffLat, dropoffLng float64) (*domain.Ride, error) {
	// 乗客の存在確認
	_, err := u.userRepo.GetByID(passengerID)
	if err != nil {
		return nil, err
	}

	// 乗客が既に進行中の配車を持っていないかチェック
	activeRide, err := u.rideRepo.GetActiveByPassengerID(passengerID)
	if err != nil && err != domain.ErrRideNotFound {
		return nil, err
	}
	if activeRide != nil {
		return nil, domain.ErrActiveRideExists
	}

	// 配車リクエストを作成
	ride, err := domain.NewRide(passengerID, pickupLat, pickupLng, dropoffLat, dropoffLng)
	if err != nil {
		return nil, err
	}

	// リポジトリ経由でデータベースに保存
	if err := u.rideRepo.Create(ride); err != nil {
		return nil, err
	}

	return ride, nil
}

// GetRide は指定されたIDの配車を取得します
func (u *rideUsecase) GetRide(id string) (*domain.Ride, error) {
	if id == "" {
		return nil, domain.ErrInvalidInput
	}
	return u.rideRepo.GetByID(id)
}

// AcceptRide はドライバーが配車を承諾します
func (u *rideUsecase) AcceptRide(rideID, driverID string) (*domain.Ride, error) {
	if rideID == "" || driverID == "" {
		return nil, domain.ErrInvalidInput
	}

	// 配車を取得
	ride, err := u.rideRepo.GetByID(rideID)
	if err != nil {
		return nil, err
	}

	// ドライバーの存在と状態を確認
	driver, err := u.driverRepo.GetByID(driverID)
	if err != nil {
		return nil, err
	}

	// ドライバーが配車可能かチェック
	if !driver.IsAvailable() {
		return nil, domain.ErrDriverNotAvailable
	}

	// ドライバーが既に進行中の配車を持っていないかチェック
	activeRide, err := u.rideRepo.GetActiveByDriverID(driverID)
	if err != nil && err != domain.ErrRideNotFound {
		return nil, err
	}
	if activeRide != nil {
		return nil, domain.ErrActiveRideExists
	}

	// 配車を承諾
	if err := ride.Accept(driverID); err != nil {
		return nil, err
	}

	// ドライバーのステータスを busy に更新
	if err := driver.UpdateStatus(domain.DriverStatusBusy); err != nil {
		return nil, err
	}

	// トランザクション内で更新を保存
	err = u.txManager.ExecuteInTransaction(func() error {
		if err := u.rideRepo.Update(ride); err != nil {
			return err
		}
		if err := u.driverRepo.Update(driver); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return ride, nil
}

// ArriveAtPickup はドライバーが乗車地点に到着したことを記録します
func (u *rideUsecase) ArriveAtPickup(rideID string) (*domain.Ride, error) {
	if rideID == "" {
		return nil, domain.ErrInvalidInput
	}

	ride, err := u.rideRepo.GetByID(rideID)
	if err != nil {
		return nil, err
	}

	if err := ride.Arrive(); err != nil {
		return nil, err
	}

	if err := u.rideRepo.Update(ride); err != nil {
		return nil, err
	}

	return ride, nil
}

// StartRide は乗車を開始します
func (u *rideUsecase) StartRide(rideID string) (*domain.Ride, error) {
	if rideID == "" {
		return nil, domain.ErrInvalidInput
	}

	ride, err := u.rideRepo.GetByID(rideID)
	if err != nil {
		return nil, err
	}

	if err := ride.Start(); err != nil {
		return nil, err
	}

	if err := u.rideRepo.Update(ride); err != nil {
		return nil, err
	}

	return ride, nil
}

// CompleteRide は乗車を完了します
func (u *rideUsecase) CompleteRide(rideID string) (*domain.Ride, error) {
	if rideID == "" {
		return nil, domain.ErrInvalidInput
	}

	ride, err := u.rideRepo.GetByID(rideID)
	if err != nil {
		return nil, err
	}

	if err := ride.Complete(); err != nil {
		return nil, err
	}

	// ドライバーのステータスを available に戻す
	var driver *domain.Driver
	if ride.DriverID != "" {
		driver, err = u.driverRepo.GetByID(ride.DriverID)
		if err != nil {
			return nil, err
		}
		if err := driver.UpdateStatus(domain.DriverStatusAvailable); err != nil {
			return nil, err
		}
	}

	// トランザクション内で更新を保存
	err = u.txManager.ExecuteInTransaction(func() error {
		if err := u.rideRepo.Update(ride); err != nil {
			return err
		}
		if driver != nil {
			if err := u.driverRepo.Update(driver); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return ride, nil
}

// CancelRide は配車をキャンセルします
func (u *rideUsecase) CancelRide(rideID string) (*domain.Ride, error) {
	if rideID == "" {
		return nil, domain.ErrInvalidInput
	}

	ride, err := u.rideRepo.GetByID(rideID)
	if err != nil {
		return nil, err
	}

	if err := ride.Cancel(); err != nil {
		return nil, err
	}

	// ドライバーが割り当てられている場合、ステータスを available に戻す
	var driver *domain.Driver
	if ride.DriverID != "" {
		driver, err = u.driverRepo.GetByID(ride.DriverID)
		if err != nil {
			return nil, err
		}
		if err := driver.UpdateStatus(domain.DriverStatusAvailable); err != nil {
			return nil, err
		}
	}

	// トランザクション内で更新を保存
	err = u.txManager.ExecuteInTransaction(func() error {
		if err := u.rideRepo.Update(ride); err != nil {
			return err
		}
		if driver != nil {
			if err := u.driverRepo.Update(driver); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return ride, nil
}

// GetPassengerRides は乗客の配車履歴を取得します
func (u *rideUsecase) GetPassengerRides(passengerID string) ([]*domain.Ride, error) {
	if passengerID == "" {
		return nil, domain.ErrInvalidInput
	}
	return u.rideRepo.GetByPassengerID(passengerID)
}

// GetDriverRides はドライバーの配車履歴を取得します
func (u *rideUsecase) GetDriverRides(driverID string) ([]*domain.Ride, error) {
	if driverID == "" {
		return nil, domain.ErrInvalidInput
	}
	return u.rideRepo.GetByDriverID(driverID)
}

// EstimateFare は乗降車地点から料金を見積もります
func (u *rideUsecase) EstimateFare(pickupLat, pickupLng, dropoffLat, dropoffLng float64) (*domain.FareBreakdown, error) {
	// 緯度経度のバリデーション
	if pickupLat < -90 || pickupLat > 90 || dropoffLat < -90 || dropoffLat > 90 {
		return nil, domain.ErrInvalidLatitude
	}
	if pickupLng < -180 || pickupLng > 180 || dropoffLng < -180 || dropoffLng > 180 {
		return nil, domain.ErrInvalidLongitude
	}

	fare := domain.CalculateRideFare(pickupLat, pickupLng, dropoffLat, dropoffLng)
	return fare, nil
}
