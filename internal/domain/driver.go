// Package domain はビジネスロジックの中核となるドメイン層を定義します
package domain

import (
	"crypto/rand"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

// DriverStatus はドライバーの稼働状態を表します
type DriverStatus string

const (
	// DriverStatusAvailable は空車（配車可能）状態
	DriverStatusAvailable DriverStatus = "available"
	// DriverStatusBusy は実車（乗客対応中）状態
	DriverStatusBusy DriverStatus = "busy"
	// DriverStatusOffline はオフライン（稼働停止）状態
	DriverStatusOffline DriverStatus = "offline"
)

// Driver はドライバーエンティティを表します
type Driver struct {
	ID            string       `json:"id" gorm:"primaryKey;type:varchar(26)"`
	Name          string       `json:"name" gorm:"type:varchar(100)"`
	Email         string       `json:"email" gorm:"type:varchar(255);uniqueIndex"`
	LicenseNumber string       `json:"license_number" gorm:"type:varchar(50);uniqueIndex"`
	Status        DriverStatus `json:"status" gorm:"type:varchar(20);default:offline"`
	Latitude      float64      `json:"latitude"`
	Longitude     float64      `json:"longitude"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

// DriverRepository はドライバーデータの永続化に関するインターフェースです
type DriverRepository interface {
	GetByID(id string) (*Driver, error)
	Create(driver *Driver) error
	Update(driver *Driver) error
	GetAll() ([]*Driver, error)
	GetAvailable() ([]*Driver, error)
	ExistsByEmail(email string) (bool, error)
	ExistsByLicenseNumber(licenseNumber string) (bool, error)
}

// NewDriver は新しいドライバーを作成します（バリデーション付き）
func NewDriver(name, email, licenseNumber string) (*Driver, error) {
	// 名前のバリデーション
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return nil, ErrEmptyName
	}

	// メールアドレスのバリデーション
	trimmedEmail := strings.TrimSpace(email)
	if !IsValidEmail(trimmedEmail) {
		return nil, ErrInvalidEmail
	}

	// 免許番号のバリデーション
	trimmedLicense := strings.TrimSpace(licenseNumber)
	if trimmedLicense == "" {
		return nil, ErrEmptyLicenseNumber
	}

	// ULID の生成
	entropy := ulid.Monotonic(rand.Reader, 0)
	id, err := ulid.New(ulid.Timestamp(time.Now()), entropy)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Driver{
		ID:            id.String(),
		Name:          trimmedName,
		Email:         trimmedEmail,
		LicenseNumber: trimmedLicense,
		Status:        DriverStatusOffline,
		Latitude:      0,
		Longitude:     0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// UpdateLocation はドライバーの現在地を更新します
func (d *Driver) UpdateLocation(latitude, longitude float64) error {
	if latitude < -90 || latitude > 90 {
		return ErrInvalidLatitude
	}
	if longitude < -180 || longitude > 180 {
		return ErrInvalidLongitude
	}

	d.Latitude = latitude
	d.Longitude = longitude
	d.UpdatedAt = time.Now()
	return nil
}

// UpdateStatus はドライバーのステータスを更新します
func (d *Driver) UpdateStatus(status DriverStatus) error {
	switch status {
	case DriverStatusAvailable, DriverStatusBusy, DriverStatusOffline:
		d.Status = status
		d.UpdatedAt = time.Now()
		return nil
	default:
		return ErrInvalidDriverStatus
	}
}

// IsAvailable はドライバーが配車可能かどうかを返します
func (d *Driver) IsAvailable() bool {
	return d.Status == DriverStatusAvailable
}
