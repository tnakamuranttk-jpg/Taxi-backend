// Package usecase はアプリケーションのビジネスロジックを実装します
package usecase

import (
	"github.com/tnakamura/taxi-backend/internal/domain"
)

// DriverUsecase はドライバーに関するビジネスロジックを定義するインターフェースです
type DriverUsecase interface {
	// GetDriver は指定されたIDのドライバーを取得します
	GetDriver(id string) (*domain.Driver, error)
	// CreateDriver は新しいドライバーを作成します
	CreateDriver(name, email, licenseNumber string) (*domain.Driver, error)
	// GetAllDrivers は全てのドライバーを取得します
	GetAllDrivers() ([]*domain.Driver, error)
	// GetAvailableDrivers は配車可能なドライバーを取得します
	GetAvailableDrivers() ([]*domain.Driver, error)
	// UpdateDriverLocation はドライバーの現在地を更新します
	UpdateDriverLocation(id string, latitude, longitude float64) (*domain.Driver, error)
	// UpdateDriverStatus はドライバーのステータスを更新します
	UpdateDriverStatus(id string, status domain.DriverStatus) (*domain.Driver, error)
}

// driverUsecase は DriverUsecase インターフェースの具体的な実装です
type driverUsecase struct {
	driverRepo domain.DriverRepository
}

// NewDriverUsecase は新しい DriverUsecase を生成します
func NewDriverUsecase(repo domain.DriverRepository) DriverUsecase {
	return &driverUsecase{
		driverRepo: repo,
	}
}

// GetDriver は指定されたIDのドライバーを取得します
func (u *driverUsecase) GetDriver(id string) (*domain.Driver, error) {
	if id == "" {
		return nil, domain.ErrInvalidInput
	}
	return u.driverRepo.GetByID(id)
}

// CreateDriver は新しいドライバーを作成して保存します
func (u *driverUsecase) CreateDriver(name, email, licenseNumber string) (*domain.Driver, error) {
	// メールアドレスの重複チェック
	exists, err := u.driverRepo.ExistsByEmail(email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrDriverAlreadyExists
	}

	// 免許番号の重複チェック
	exists, err = u.driverRepo.ExistsByLicenseNumber(licenseNumber)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrDriverAlreadyExists
	}

	// ドメイン層の NewDriver を呼び出して ULID を生成したドライバーオブジェクトを作成
	driver, err := domain.NewDriver(name, email, licenseNumber)
	if err != nil {
		return nil, err
	}

	// リポジトリ経由でデータベースに保存
	if err := u.driverRepo.Create(driver); err != nil {
		return nil, err
	}

	return driver, nil
}

// GetAllDrivers は登録されている全ドライバーを取得します
func (u *driverUsecase) GetAllDrivers() ([]*domain.Driver, error) {
	return u.driverRepo.GetAll()
}

// GetAvailableDrivers は配車可能なドライバーを取得します
func (u *driverUsecase) GetAvailableDrivers() ([]*domain.Driver, error) {
	return u.driverRepo.GetAvailable()
}

// UpdateDriverLocation はドライバーの現在地を更新します
func (u *driverUsecase) UpdateDriverLocation(id string, latitude, longitude float64) (*domain.Driver, error) {
	if id == "" {
		return nil, domain.ErrInvalidInput
	}

	driver, err := u.driverRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := driver.UpdateLocation(latitude, longitude); err != nil {
		return nil, err
	}

	if err := u.driverRepo.Update(driver); err != nil {
		return nil, err
	}

	return driver, nil
}

// UpdateDriverStatus はドライバーのステータスを更新します
func (u *driverUsecase) UpdateDriverStatus(id string, status domain.DriverStatus) (*domain.Driver, error) {
	if id == "" {
		return nil, domain.ErrInvalidInput
	}

	driver, err := u.driverRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := driver.UpdateStatus(status); err != nil {
		return nil, err
	}

	if err := u.driverRepo.Update(driver); err != nil {
		return nil, err
	}

	return driver, nil
}
