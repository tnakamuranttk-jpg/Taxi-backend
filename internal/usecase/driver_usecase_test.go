package usecase_test

import (
	"testing"

	"github.com/tnakamura/taxi-backend/internal/domain"
	"github.com/tnakamura/taxi-backend/internal/usecase"
)

// MockDriverRepository はテスト用のモックリポジトリです
type MockDriverRepository struct {
	drivers     map[string]*domain.Driver
	shouldError bool
}

func NewMockDriverRepository() *MockDriverRepository {
	return &MockDriverRepository{
		drivers: make(map[string]*domain.Driver),
	}
}

func (m *MockDriverRepository) GetByID(id string) (*domain.Driver, error) {
	if m.shouldError {
		return nil, domain.ErrDriverNotFound
	}
	driver, ok := m.drivers[id]
	if !ok {
		return nil, domain.ErrDriverNotFound
	}
	return driver, nil
}

func (m *MockDriverRepository) Create(driver *domain.Driver) error {
	if m.shouldError {
		return domain.ErrDriverAlreadyExists
	}
	m.drivers[driver.ID] = driver
	return nil
}

func (m *MockDriverRepository) Update(driver *domain.Driver) error {
	if m.shouldError {
		return domain.ErrDriverNotFound
	}
	if _, ok := m.drivers[driver.ID]; !ok {
		return domain.ErrDriverNotFound
	}
	m.drivers[driver.ID] = driver
	return nil
}

func (m *MockDriverRepository) GetAll() ([]*domain.Driver, error) {
	if m.shouldError {
		return nil, domain.ErrInvalidInput
	}
	result := make([]*domain.Driver, 0, len(m.drivers))
	for _, driver := range m.drivers {
		result = append(result, driver)
	}
	return result, nil
}

func (m *MockDriverRepository) GetAvailable() ([]*domain.Driver, error) {
	if m.shouldError {
		return nil, domain.ErrInvalidInput
	}
	result := make([]*domain.Driver, 0)
	for _, driver := range m.drivers {
		if driver.Status == domain.DriverStatusAvailable {
			result = append(result, driver)
		}
	}
	return result, nil
}

func (m *MockDriverRepository) ExistsByEmail(email string) (bool, error) {
	for _, driver := range m.drivers {
		if driver.Email == email {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockDriverRepository) ExistsByLicenseNumber(licenseNumber string) (bool, error) {
	for _, driver := range m.drivers {
		if driver.LicenseNumber == licenseNumber {
			return true, nil
		}
	}
	return false, nil
}

// === ドライバー作成テスト ===

func TestCreateDriver_Success(t *testing.T) {
	mockRepo := NewMockDriverRepository()
	uc := usecase.NewDriverUsecase(mockRepo)

	driver, err := uc.CreateDriver("佐藤運転手", "sato@example.com", "LICENSE-001")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if driver == nil {
		t.Fatal("expected driver to be created")
	}

	if driver.ID == "" {
		t.Error("expected driver ID to be generated (ULID)")
	}

	if driver.Name != "佐藤運転手" {
		t.Errorf("expected name '佐藤運転手', got '%s'", driver.Name)
	}

	if driver.Email != "sato@example.com" {
		t.Errorf("expected email 'sato@example.com', got '%s'", driver.Email)
	}

	if driver.LicenseNumber != "LICENSE-001" {
		t.Errorf("expected license number 'LICENSE-001', got '%s'", driver.LicenseNumber)
	}

	if driver.Status != domain.DriverStatusOffline {
		t.Errorf("expected status 'offline', got '%s'", driver.Status)
	}
}

func TestCreateDriver_EmptyName(t *testing.T) {
	mockRepo := NewMockDriverRepository()
	uc := usecase.NewDriverUsecase(mockRepo)

	_, err := uc.CreateDriver("", "sato@example.com", "LICENSE-001")

	if err != domain.ErrEmptyName {
		t.Errorf("expected ErrEmptyName, got %v", err)
	}
}

func TestCreateDriver_InvalidEmail(t *testing.T) {
	mockRepo := NewMockDriverRepository()
	uc := usecase.NewDriverUsecase(mockRepo)

	_, err := uc.CreateDriver("佐藤運転手", "invalid-email", "LICENSE-001")

	if err != domain.ErrInvalidEmail {
		t.Errorf("expected ErrInvalidEmail, got %v", err)
	}
}

func TestCreateDriver_EmptyLicenseNumber(t *testing.T) {
	mockRepo := NewMockDriverRepository()
	uc := usecase.NewDriverUsecase(mockRepo)

	_, err := uc.CreateDriver("佐藤運転手", "sato@example.com", "")

	if err != domain.ErrEmptyLicenseNumber {
		t.Errorf("expected ErrEmptyLicenseNumber, got %v", err)
	}
}

func TestCreateDriver_DuplicateEmail(t *testing.T) {
	mockRepo := NewMockDriverRepository()
	uc := usecase.NewDriverUsecase(mockRepo)

	// 最初のドライバーを作成
	_, err := uc.CreateDriver("佐藤運転手", "sato@example.com", "LICENSE-001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 同じメールアドレスで2人目を作成しようとする
	_, err = uc.CreateDriver("田中運転手", "sato@example.com", "LICENSE-002")
	if err != domain.ErrDriverAlreadyExists {
		t.Errorf("expected ErrDriverAlreadyExists, got %v", err)
	}
}

func TestCreateDriver_DuplicateLicenseNumber(t *testing.T) {
	mockRepo := NewMockDriverRepository()
	uc := usecase.NewDriverUsecase(mockRepo)

	// 最初のドライバーを作成
	_, err := uc.CreateDriver("佐藤運転手", "sato@example.com", "LICENSE-001")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 同じ免許番号で2人目を作成しようとする
	_, err = uc.CreateDriver("田中運転手", "tanaka@example.com", "LICENSE-001")
	if err != domain.ErrDriverAlreadyExists {
		t.Errorf("expected ErrDriverAlreadyExists, got %v", err)
	}
}

// === ドライバー取得テスト ===

func TestGetDriver_Success(t *testing.T) {
	mockRepo := NewMockDriverRepository()
	uc := usecase.NewDriverUsecase(mockRepo)

	// ドライバーを作成
	createdDriver, _ := uc.CreateDriver("佐藤運転手", "sato@example.com", "LICENSE-001")

	// 取得
	driver, err := uc.GetDriver(createdDriver.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if driver.ID != createdDriver.ID {
		t.Errorf("expected ID '%s', got '%s'", createdDriver.ID, driver.ID)
	}
}

func TestGetDriver_NotFound(t *testing.T) {
	mockRepo := NewMockDriverRepository()
	uc := usecase.NewDriverUsecase(mockRepo)

	_, err := uc.GetDriver("nonexistent-id")

	if err != domain.ErrDriverNotFound {
		t.Errorf("expected ErrDriverNotFound, got %v", err)
	}
}

func TestGetDriver_EmptyID(t *testing.T) {
	mockRepo := NewMockDriverRepository()
	uc := usecase.NewDriverUsecase(mockRepo)

	_, err := uc.GetDriver("")

	if err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

// === ドライバー一覧取得テスト ===

func TestGetAllDrivers_Success(t *testing.T) {
	mockRepo := NewMockDriverRepository()
	uc := usecase.NewDriverUsecase(mockRepo)

	// ドライバーを複数作成
	uc.CreateDriver("ドライバー1", "driver1@example.com", "LICENSE-001")
	uc.CreateDriver("ドライバー2", "driver2@example.com", "LICENSE-002")

	drivers, err := uc.GetAllDrivers()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(drivers) != 2 {
		t.Errorf("expected 2 drivers, got %d", len(drivers))
	}
}

func TestGetAllDrivers_Empty(t *testing.T) {
	mockRepo := NewMockDriverRepository()
	uc := usecase.NewDriverUsecase(mockRepo)

	drivers, err := uc.GetAllDrivers()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(drivers) != 0 {
		t.Errorf("expected 0 drivers, got %d", len(drivers))
	}
}

// === 配車可能ドライバー取得テスト ===

func TestGetAvailableDrivers_Success(t *testing.T) {
	mockRepo := NewMockDriverRepository()
	uc := usecase.NewDriverUsecase(mockRepo)

	// ドライバーを作成
	driver1, _ := uc.CreateDriver("ドライバー1", "driver1@example.com", "LICENSE-001")
	driver2, _ := uc.CreateDriver("ドライバー2", "driver2@example.com", "LICENSE-002")

	// ステータスを更新
	uc.UpdateDriverStatus(driver1.ID, domain.DriverStatusAvailable)
	uc.UpdateDriverStatus(driver2.ID, domain.DriverStatusBusy)

	// 配車可能なドライバーを取得
	availableDrivers, err := uc.GetAvailableDrivers()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(availableDrivers) != 1 {
		t.Errorf("expected 1 available driver, got %d", len(availableDrivers))
	}
}

// === 位置更新テスト ===

func TestUpdateDriverLocation_Success(t *testing.T) {
	mockRepo := NewMockDriverRepository()
	uc := usecase.NewDriverUsecase(mockRepo)

	// ドライバーを作成
	createdDriver, _ := uc.CreateDriver("佐藤運転手", "sato@example.com", "LICENSE-001")

	// 位置を更新
	driver, err := uc.UpdateDriverLocation(createdDriver.ID, 35.6762, 139.6503)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if driver.Latitude != 35.6762 {
		t.Errorf("expected latitude 35.6762, got %f", driver.Latitude)
	}

	if driver.Longitude != 139.6503 {
		t.Errorf("expected longitude 139.6503, got %f", driver.Longitude)
	}
}

func TestUpdateDriverLocation_InvalidLatitude(t *testing.T) {
	mockRepo := NewMockDriverRepository()
	uc := usecase.NewDriverUsecase(mockRepo)

	// ドライバーを作成
	createdDriver, _ := uc.CreateDriver("佐藤運転手", "sato@example.com", "LICENSE-001")

	// 無効な緯度で更新
	_, err := uc.UpdateDriverLocation(createdDriver.ID, 100.0, 139.6503)

	if err != domain.ErrInvalidLatitude {
		t.Errorf("expected ErrInvalidLatitude, got %v", err)
	}
}

func TestUpdateDriverLocation_InvalidLongitude(t *testing.T) {
	mockRepo := NewMockDriverRepository()
	uc := usecase.NewDriverUsecase(mockRepo)

	// ドライバーを作成
	createdDriver, _ := uc.CreateDriver("佐藤運転手", "sato@example.com", "LICENSE-001")

	// 無効な経度で更新
	_, err := uc.UpdateDriverLocation(createdDriver.ID, 35.6762, 200.0)

	if err != domain.ErrInvalidLongitude {
		t.Errorf("expected ErrInvalidLongitude, got %v", err)
	}
}

// === ステータス更新テスト ===

func TestUpdateDriverStatus_Success(t *testing.T) {
	mockRepo := NewMockDriverRepository()
	uc := usecase.NewDriverUsecase(mockRepo)

	// ドライバーを作成
	createdDriver, _ := uc.CreateDriver("佐藤運転手", "sato@example.com", "LICENSE-001")

	// ステータスを更新
	driver, err := uc.UpdateDriverStatus(createdDriver.ID, domain.DriverStatusAvailable)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if driver.Status != domain.DriverStatusAvailable {
		t.Errorf("expected status 'available', got '%s'", driver.Status)
	}
}

func TestUpdateDriverStatus_InvalidStatus(t *testing.T) {
	mockRepo := NewMockDriverRepository()
	uc := usecase.NewDriverUsecase(mockRepo)

	// ドライバーを作成
	createdDriver, _ := uc.CreateDriver("佐藤運転手", "sato@example.com", "LICENSE-001")

	// 無効なステータスで更新
	_, err := uc.UpdateDriverStatus(createdDriver.ID, "invalid-status")

	if err != domain.ErrInvalidDriverStatus {
		t.Errorf("expected ErrInvalidDriverStatus, got %v", err)
	}
}

func TestUpdateDriverStatus_NotFound(t *testing.T) {
	mockRepo := NewMockDriverRepository()
	uc := usecase.NewDriverUsecase(mockRepo)

	_, err := uc.UpdateDriverStatus("nonexistent-id", domain.DriverStatusAvailable)

	if err != domain.ErrDriverNotFound {
		t.Errorf("expected ErrDriverNotFound, got %v", err)
	}
}
