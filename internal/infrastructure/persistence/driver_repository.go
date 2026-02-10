// Package persistence はデータベースとの永続化処理を実装します
package persistence

import (
	"errors"

	"github.com/tnakamura/taxi-backend/internal/domain"
	"gorm.io/gorm"
)

// driverRepository は domain.DriverRepository の具体的な実装です
type driverRepository struct {
	db *gorm.DB
}

// NewDriverRepository は新しい DriverRepository を生成します
func NewDriverRepository(db *gorm.DB) domain.DriverRepository {
	return &driverRepository{
		db: db,
	}
}

// GetByID は指定されたIDのドライバーを取得します
func (r *driverRepository) GetByID(id string) (*domain.Driver, error) {
	var driver domain.Driver
	if err := r.db.First(&driver, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrDriverNotFound
		}
		return nil, err
	}
	return &driver, nil
}

// Create は新しいドライバーをデータベースに保存します
func (r *driverRepository) Create(driver *domain.Driver) error {
	if err := r.db.Create(driver).Error; err != nil {
		if isDuplicateError(err) {
			return domain.ErrDriverAlreadyExists
		}
		return err
	}
	return nil
}

// Update はドライバー情報を更新します
func (r *driverRepository) Update(driver *domain.Driver) error {
	result := r.db.Save(driver)
	if result.Error != nil {
		if isDuplicateError(result.Error) {
			return domain.ErrDriverAlreadyExists
		}
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrDriverNotFound
	}
	return nil
}

// GetAll は全てのドライバーを取得します
func (r *driverRepository) GetAll() ([]*domain.Driver, error) {
	var drivers []*domain.Driver
	if err := r.db.Find(&drivers).Error; err != nil {
		return nil, err
	}
	return drivers, nil
}

// GetAvailable は配車可能なドライバーを取得します
func (r *driverRepository) GetAvailable() ([]*domain.Driver, error) {
	var drivers []*domain.Driver
	if err := r.db.Where("status = ?", domain.DriverStatusAvailable).Find(&drivers).Error; err != nil {
		return nil, err
	}
	return drivers, nil
}

// ExistsByEmail は指定されたメールアドレスのドライバーが存在するかチェックします
func (r *driverRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	if err := r.db.Model(&domain.Driver{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByLicenseNumber は指定された免許番号のドライバーが存在するかチェックします
func (r *driverRepository) ExistsByLicenseNumber(licenseNumber string) (bool, error) {
	var count int64
	if err := r.db.Model(&domain.Driver{}).Where("license_number = ?", licenseNumber).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
