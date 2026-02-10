// Package usecase はアプリケーションのビジネスロジックを実装します
package usecase

import (
	"testing"
	"time"

	"github.com/tnakamura/taxi-backend/internal/domain"
)

// ========== Ride モックリポジトリ ==========

type mockRideRepository struct {
	rides          map[string]*domain.Ride
	passengerRides map[string][]*domain.Ride
	driverRides    map[string][]*domain.Ride
}

func newMockRideRepository() *mockRideRepository {
	return &mockRideRepository{
		rides:          make(map[string]*domain.Ride),
		passengerRides: make(map[string][]*domain.Ride),
		driverRides:    make(map[string][]*domain.Ride),
	}
}

func (m *mockRideRepository) GetByID(id string) (*domain.Ride, error) {
	if ride, ok := m.rides[id]; ok {
		return ride, nil
	}
	return nil, domain.ErrRideNotFound
}

func (m *mockRideRepository) Create(ride *domain.Ride) error {
	m.rides[ride.ID] = ride
	m.passengerRides[ride.PassengerID] = append(m.passengerRides[ride.PassengerID], ride)
	return nil
}

func (m *mockRideRepository) Update(ride *domain.Ride) error {
	if _, ok := m.rides[ride.ID]; !ok {
		return domain.ErrRideNotFound
	}
	m.rides[ride.ID] = ride
	return nil
}

func (m *mockRideRepository) GetByPassengerID(passengerID string) ([]*domain.Ride, error) {
	return m.passengerRides[passengerID], nil
}

func (m *mockRideRepository) GetByDriverID(driverID string) ([]*domain.Ride, error) {
	return m.driverRides[driverID], nil
}

func (m *mockRideRepository) GetActiveByPassengerID(passengerID string) (*domain.Ride, error) {
	for _, ride := range m.passengerRides[passengerID] {
		if ride.IsActive() {
			return ride, nil
		}
	}
	return nil, domain.ErrRideNotFound
}

func (m *mockRideRepository) GetActiveByDriverID(driverID string) (*domain.Ride, error) {
	for _, ride := range m.rides {
		if ride.DriverID == driverID && ride.IsActive() {
			return ride, nil
		}
	}
	return nil, domain.ErrRideNotFound
}

// ========== User モックリポジトリ（Ride用） ==========

type mockUserRepoForRide struct {
	users map[string]*domain.User
}

func newMockUserRepoForRide() *mockUserRepoForRide {
	return &mockUserRepoForRide{
		users: make(map[string]*domain.User),
	}
}

func (m *mockUserRepoForRide) GetByID(id string) (*domain.User, error) {
	if user, ok := m.users[id]; ok {
		return user, nil
	}
	return nil, domain.ErrUserNotFound
}

func (m *mockUserRepoForRide) Create(user *domain.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepoForRide) GetAll() ([]*domain.User, error) {
	var result []*domain.User
	for _, u := range m.users {
		result = append(result, u)
	}
	return result, nil
}

func (m *mockUserRepoForRide) ExistsByEmail(email string) (bool, error) {
	for _, u := range m.users {
		if u.Email == email {
			return true, nil
		}
	}
	return false, nil
}

// ========== Driver モックリポジトリ（Ride用） ==========

type mockDriverRepoForRide struct {
	drivers map[string]*domain.Driver
}

func newMockDriverRepoForRide() *mockDriverRepoForRide {
	return &mockDriverRepoForRide{
		drivers: make(map[string]*domain.Driver),
	}
}

func (m *mockDriverRepoForRide) GetByID(id string) (*domain.Driver, error) {
	if driver, ok := m.drivers[id]; ok {
		return driver, nil
	}
	return nil, domain.ErrDriverNotFound
}

func (m *mockDriverRepoForRide) Create(driver *domain.Driver) error {
	m.drivers[driver.ID] = driver
	return nil
}

func (m *mockDriverRepoForRide) Update(driver *domain.Driver) error {
	if _, ok := m.drivers[driver.ID]; !ok {
		return domain.ErrDriverNotFound
	}
	m.drivers[driver.ID] = driver
	return nil
}

func (m *mockDriverRepoForRide) GetAll() ([]*domain.Driver, error) {
	var result []*domain.Driver
	for _, d := range m.drivers {
		result = append(result, d)
	}
	return result, nil
}

func (m *mockDriverRepoForRide) GetAvailable() ([]*domain.Driver, error) {
	var result []*domain.Driver
	for _, d := range m.drivers {
		if d.IsAvailable() {
			result = append(result, d)
		}
	}
	return result, nil
}

func (m *mockDriverRepoForRide) ExistsByEmail(email string) (bool, error) {
	return false, nil
}

func (m *mockDriverRepoForRide) ExistsByLicenseNumber(licenseNumber string) (bool, error) {
	return false, nil
}

// ========== モックTransactionManager ==========

type mockTransactionManager struct{}

func (m *mockTransactionManager) ExecuteInTransaction(fn func() error) error {
	return fn()
}

// ========== テストヘルパー ==========

func createTestUser() *domain.User {
	return &domain.User{
		ID:    "test-user-id",
		Name:  "Test User",
		Email: "test@example.com",
	}
}

func createTestDriver() *domain.Driver {
	return &domain.Driver{
		ID:            "test-driver-id",
		Name:          "Test Driver",
		Email:         "driver@example.com",
		LicenseNumber: "LICENSE-001",
		Status:        domain.DriverStatusAvailable,
		Latitude:      35.6762,
		Longitude:     139.6503,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// ========== テスト ==========

// --- CreateRide テスト ---

func TestCreateRide_Success(t *testing.T) {
	rideRepo := newMockRideRepository()
	userRepo := newMockUserRepoForRide()
	driverRepo := newMockDriverRepoForRide()
	txManager := &mockTransactionManager{}

	user := createTestUser()
	userRepo.users[user.ID] = user

	usecase := NewRideUsecase(rideRepo, driverRepo, userRepo, txManager)

	ride, err := usecase.CreateRide(user.ID, 35.6762, 139.6503, 35.6895, 139.6917)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ride == nil {
		t.Fatal("expected ride, got nil")
	}
	if ride.PassengerID != user.ID {
		t.Errorf("expected passenger ID %s, got %s", user.ID, ride.PassengerID)
	}
	if ride.Status != domain.RideStatusRequested {
		t.Errorf("expected status %s, got %s", domain.RideStatusRequested, ride.Status)
	}
}

func TestCreateRide_UserNotFound(t *testing.T) {
	rideRepo := newMockRideRepository()
	userRepo := newMockUserRepoForRide()
	driverRepo := newMockDriverRepoForRide()
	txManager := &mockTransactionManager{}

	usecase := NewRideUsecase(rideRepo, driverRepo, userRepo, txManager)

	_, err := usecase.CreateRide("nonexistent-user", 35.6762, 139.6503, 35.6895, 139.6917)
	if err != domain.ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestCreateRide_ActiveRideExists(t *testing.T) {
	rideRepo := newMockRideRepository()
	userRepo := newMockUserRepoForRide()
	driverRepo := newMockDriverRepoForRide()
	txManager := &mockTransactionManager{}

	user := createTestUser()
	userRepo.users[user.ID] = user

	// 既存の進行中配車を作成
	existingRide := &domain.Ride{
		ID:          "existing-ride-id",
		PassengerID: user.ID,
		Status:      domain.RideStatusRequested,
	}
	rideRepo.rides[existingRide.ID] = existingRide
	rideRepo.passengerRides[user.ID] = append(rideRepo.passengerRides[user.ID], existingRide)

	usecase := NewRideUsecase(rideRepo, driverRepo, userRepo, txManager)

	_, err := usecase.CreateRide(user.ID, 35.6762, 139.6503, 35.6895, 139.6917)
	if err != domain.ErrActiveRideExists {
		t.Errorf("expected ErrActiveRideExists, got %v", err)
	}
}

func TestCreateRide_InvalidLatitude(t *testing.T) {
	rideRepo := newMockRideRepository()
	userRepo := newMockUserRepoForRide()
	driverRepo := newMockDriverRepoForRide()
	txManager := &mockTransactionManager{}

	user := createTestUser()
	userRepo.users[user.ID] = user

	usecase := NewRideUsecase(rideRepo, driverRepo, userRepo, txManager)

	_, err := usecase.CreateRide(user.ID, 100, 139.6503, 35.6895, 139.6917)
	if err != domain.ErrInvalidLatitude {
		t.Errorf("expected ErrInvalidLatitude, got %v", err)
	}
}

// --- AcceptRide テスト ---

func TestAcceptRide_Success(t *testing.T) {
	rideRepo := newMockRideRepository()
	userRepo := newMockUserRepoForRide()
	driverRepo := newMockDriverRepoForRide()
	txManager := &mockTransactionManager{}

	user := createTestUser()
	userRepo.users[user.ID] = user

	driver := createTestDriver()
	driverRepo.drivers[driver.ID] = driver

	ride := &domain.Ride{
		ID:          "test-ride-id",
		PassengerID: user.ID,
		Status:      domain.RideStatusRequested,
	}
	rideRepo.rides[ride.ID] = ride

	usecase := NewRideUsecase(rideRepo, driverRepo, userRepo, txManager)

	acceptedRide, err := usecase.AcceptRide(ride.ID, driver.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if acceptedRide.Status != domain.RideStatusAccepted {
		t.Errorf("expected status %s, got %s", domain.RideStatusAccepted, acceptedRide.Status)
	}
	if acceptedRide.DriverID != driver.ID {
		t.Errorf("expected driver ID %s, got %s", driver.ID, acceptedRide.DriverID)
	}
	if driverRepo.drivers[driver.ID].Status != domain.DriverStatusBusy {
		t.Errorf("expected driver status %s, got %s", domain.DriverStatusBusy, driverRepo.drivers[driver.ID].Status)
	}
}

func TestAcceptRide_DriverNotAvailable(t *testing.T) {
	rideRepo := newMockRideRepository()
	userRepo := newMockUserRepoForRide()
	driverRepo := newMockDriverRepoForRide()
	txManager := &mockTransactionManager{}

	user := createTestUser()
	userRepo.users[user.ID] = user

	driver := createTestDriver()
	driver.Status = domain.DriverStatusBusy
	driverRepo.drivers[driver.ID] = driver

	ride := &domain.Ride{
		ID:          "test-ride-id",
		PassengerID: user.ID,
		Status:      domain.RideStatusRequested,
	}
	rideRepo.rides[ride.ID] = ride

	usecase := NewRideUsecase(rideRepo, driverRepo, userRepo, txManager)

	_, err := usecase.AcceptRide(ride.ID, driver.ID)
	if err != domain.ErrDriverNotAvailable {
		t.Errorf("expected ErrDriverNotAvailable, got %v", err)
	}
}

func TestAcceptRide_RideNotFound(t *testing.T) {
	rideRepo := newMockRideRepository()
	userRepo := newMockUserRepoForRide()
	driverRepo := newMockDriverRepoForRide()
	txManager := &mockTransactionManager{}

	driver := createTestDriver()
	driverRepo.drivers[driver.ID] = driver

	usecase := NewRideUsecase(rideRepo, driverRepo, userRepo, txManager)

	_, err := usecase.AcceptRide("nonexistent-ride", driver.ID)
	if err != domain.ErrRideNotFound {
		t.Errorf("expected ErrRideNotFound, got %v", err)
	}
}

// --- CompleteRide テスト ---

func TestCompleteRide_Success(t *testing.T) {
	rideRepo := newMockRideRepository()
	userRepo := newMockUserRepoForRide()
	driverRepo := newMockDriverRepoForRide()
	txManager := &mockTransactionManager{}

	driver := createTestDriver()
	driver.Status = domain.DriverStatusBusy
	driverRepo.drivers[driver.ID] = driver

	ride := &domain.Ride{
		ID:       "test-ride-id",
		DriverID: driver.ID,
		Status:   domain.RideStatusOngoing,
	}
	rideRepo.rides[ride.ID] = ride

	usecase := NewRideUsecase(rideRepo, driverRepo, userRepo, txManager)

	completedRide, err := usecase.CompleteRide(ride.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if completedRide.Status != domain.RideStatusCompleted {
		t.Errorf("expected status %s, got %s", domain.RideStatusCompleted, completedRide.Status)
	}
	if driverRepo.drivers[driver.ID].Status != domain.DriverStatusAvailable {
		t.Errorf("expected driver status %s, got %s", domain.DriverStatusAvailable, driverRepo.drivers[driver.ID].Status)
	}
}

func TestCompleteRide_InvalidStatus(t *testing.T) {
	rideRepo := newMockRideRepository()
	userRepo := newMockUserRepoForRide()
	driverRepo := newMockDriverRepoForRide()
	txManager := &mockTransactionManager{}

	ride := &domain.Ride{
		ID:     "test-ride-id",
		Status: domain.RideStatusRequested, // 完了できない状態
	}
	rideRepo.rides[ride.ID] = ride

	usecase := NewRideUsecase(rideRepo, driverRepo, userRepo, txManager)

	_, err := usecase.CompleteRide(ride.ID)
	if err != domain.ErrInvalidRideStatus {
		t.Errorf("expected ErrInvalidRideStatus, got %v", err)
	}
}

// --- CancelRide テスト ---

func TestCancelRide_Success(t *testing.T) {
	rideRepo := newMockRideRepository()
	userRepo := newMockUserRepoForRide()
	driverRepo := newMockDriverRepoForRide()
	txManager := &mockTransactionManager{}

	driver := createTestDriver()
	driver.Status = domain.DriverStatusBusy
	driverRepo.drivers[driver.ID] = driver

	ride := &domain.Ride{
		ID:       "test-ride-id",
		DriverID: driver.ID,
		Status:   domain.RideStatusAccepted,
	}
	rideRepo.rides[ride.ID] = ride

	usecase := NewRideUsecase(rideRepo, driverRepo, userRepo, txManager)

	cancelledRide, err := usecase.CancelRide(ride.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cancelledRide.Status != domain.RideStatusCancelled {
		t.Errorf("expected status %s, got %s", domain.RideStatusCancelled, cancelledRide.Status)
	}
	if driverRepo.drivers[driver.ID].Status != domain.DriverStatusAvailable {
		t.Errorf("expected driver status %s, got %s", domain.DriverStatusAvailable, driverRepo.drivers[driver.ID].Status)
	}
}

func TestCancelRide_AlreadyCompleted(t *testing.T) {
	rideRepo := newMockRideRepository()
	userRepo := newMockUserRepoForRide()
	driverRepo := newMockDriverRepoForRide()
	txManager := &mockTransactionManager{}

	ride := &domain.Ride{
		ID:     "test-ride-id",
		Status: domain.RideStatusCompleted,
	}
	rideRepo.rides[ride.ID] = ride

	usecase := NewRideUsecase(rideRepo, driverRepo, userRepo, txManager)

	_, err := usecase.CancelRide(ride.ID)
	if err != domain.ErrInvalidRideStatus {
		t.Errorf("expected ErrInvalidRideStatus, got %v", err)
	}
}

// --- GetRide テスト ---

func TestGetRide_Success(t *testing.T) {
	rideRepo := newMockRideRepository()
	userRepo := newMockUserRepoForRide()
	driverRepo := newMockDriverRepoForRide()
	txManager := &mockTransactionManager{}

	ride := &domain.Ride{
		ID:     "test-ride-id",
		Status: domain.RideStatusRequested,
	}
	rideRepo.rides[ride.ID] = ride

	usecase := NewRideUsecase(rideRepo, driverRepo, userRepo, txManager)

	foundRide, err := usecase.GetRide(ride.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if foundRide.ID != ride.ID {
		t.Errorf("expected ride ID %s, got %s", ride.ID, foundRide.ID)
	}
}

func TestGetRide_NotFound(t *testing.T) {
	rideRepo := newMockRideRepository()
	userRepo := newMockUserRepoForRide()
	driverRepo := newMockDriverRepoForRide()
	txManager := &mockTransactionManager{}

	usecase := NewRideUsecase(rideRepo, driverRepo, userRepo, txManager)

	_, err := usecase.GetRide("nonexistent-ride")
	if err != domain.ErrRideNotFound {
		t.Errorf("expected ErrRideNotFound, got %v", err)
	}
}

func TestGetRide_EmptyID(t *testing.T) {
	rideRepo := newMockRideRepository()
	userRepo := newMockUserRepoForRide()
	driverRepo := newMockDriverRepoForRide()
	txManager := &mockTransactionManager{}

	usecase := NewRideUsecase(rideRepo, driverRepo, userRepo, txManager)

	_, err := usecase.GetRide("")
	if err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

// --- GetPassengerRides テスト ---

func TestGetPassengerRides_Success(t *testing.T) {
	rideRepo := newMockRideRepository()
	userRepo := newMockUserRepoForRide()
	driverRepo := newMockDriverRepoForRide()
	txManager := &mockTransactionManager{}

	passengerID := "test-passenger-id"
	ride1 := &domain.Ride{ID: "ride-1", PassengerID: passengerID, Status: domain.RideStatusCompleted}
	ride2 := &domain.Ride{ID: "ride-2", PassengerID: passengerID, Status: domain.RideStatusRequested}
	rideRepo.passengerRides[passengerID] = []*domain.Ride{ride1, ride2}

	usecase := NewRideUsecase(rideRepo, driverRepo, userRepo, txManager)

	rides, err := usecase.GetPassengerRides(passengerID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(rides) != 2 {
		t.Errorf("expected 2 rides, got %d", len(rides))
	}
}

func TestGetPassengerRides_EmptyID(t *testing.T) {
	rideRepo := newMockRideRepository()
	userRepo := newMockUserRepoForRide()
	driverRepo := newMockDriverRepoForRide()
	txManager := &mockTransactionManager{}

	usecase := NewRideUsecase(rideRepo, driverRepo, userRepo, txManager)

	_, err := usecase.GetPassengerRides("")
	if err != domain.ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

// --- CompleteRide 料金計算テスト ---

func TestCompleteRide_FareCalculation(t *testing.T) {
	rideRepo := newMockRideRepository()
	userRepo := newMockUserRepoForRide()
	driverRepo := newMockDriverRepoForRide()
	txManager := &mockTransactionManager{}

	driver := createTestDriver()
	driver.Status = domain.DriverStatusBusy
	driverRepo.drivers[driver.ID] = driver

	// 東京駅 → 渋谷駅（約5km）
	ride := &domain.Ride{
		ID:               "test-ride-id",
		DriverID:         driver.ID,
		PickupLatitude:   35.6812,
		PickupLongitude:  139.7671,
		DropoffLatitude:  35.6580,
		DropoffLongitude: 139.7016,
		Status:           domain.RideStatusOngoing,
	}
	rideRepo.rides[ride.ID] = ride

	usecase := NewRideUsecase(rideRepo, driverRepo, userRepo, txManager)

	completedRide, err := usecase.CompleteRide(ride.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 料金が自動計算されていること
	if completedRide.FareAmount <= 0 {
		t.Error("expected positive fare amount after completion")
	}
	if completedRide.DistanceKm <= 0 {
		t.Error("expected positive distance after completion")
	}

	// 東京→渋谷は約5km、料金は1500〜2500円程度
	if completedRide.FareAmount < 1000 || completedRide.FareAmount > 3000 {
		t.Errorf("fare seems unreasonable: %d yen for %.2f km", completedRide.FareAmount, completedRide.DistanceKm)
	}
}

// --- EstimateFare テスト ---

func TestEstimateFare_Success(t *testing.T) {
	rideRepo := newMockRideRepository()
	userRepo := newMockUserRepoForRide()
	driverRepo := newMockDriverRepoForRide()
	txManager := &mockTransactionManager{}

	usecase := NewRideUsecase(rideRepo, driverRepo, userRepo, txManager)

	// 東京駅 → 渋谷駅
	fare, err := usecase.EstimateFare(35.6812, 139.7671, 35.6580, 139.7016)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if fare == nil {
		t.Fatal("expected fare breakdown, got nil")
	}
	if fare.TotalFare <= 0 {
		t.Error("expected positive total fare")
	}
	if fare.BaseFare != domain.BaseFare {
		t.Errorf("expected base fare %d, got %d", domain.BaseFare, fare.BaseFare)
	}
	if fare.DistanceKm <= 0 {
		t.Error("expected positive distance")
	}
}

func TestEstimateFare_ShortDistance(t *testing.T) {
	rideRepo := newMockRideRepository()
	userRepo := newMockUserRepoForRide()
	driverRepo := newMockDriverRepoForRide()
	txManager := &mockTransactionManager{}

	usecase := NewRideUsecase(rideRepo, driverRepo, userRepo, txManager)

	// ほぼ同一地点（距離ゼロに近い）
	fare, err := usecase.EstimateFare(35.6812, 139.7671, 35.6813, 139.7672)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if fare.TotalFare != domain.MinimumFare {
		t.Errorf("expected minimum fare %d for short distance, got %d", domain.MinimumFare, fare.TotalFare)
	}
}

func TestEstimateFare_InvalidLatitude(t *testing.T) {
	rideRepo := newMockRideRepository()
	userRepo := newMockUserRepoForRide()
	driverRepo := newMockDriverRepoForRide()
	txManager := &mockTransactionManager{}

	usecase := NewRideUsecase(rideRepo, driverRepo, userRepo, txManager)

	_, err := usecase.EstimateFare(100, 139.7671, 35.6580, 139.7016)
	if err != domain.ErrInvalidLatitude {
		t.Errorf("expected ErrInvalidLatitude, got %v", err)
	}
}

func TestEstimateFare_InvalidLongitude(t *testing.T) {
	rideRepo := newMockRideRepository()
	userRepo := newMockUserRepoForRide()
	driverRepo := newMockDriverRepoForRide()
	txManager := &mockTransactionManager{}

	usecase := NewRideUsecase(rideRepo, driverRepo, userRepo, txManager)

	_, err := usecase.EstimateFare(35.6812, 200, 35.6580, 139.7016)
	if err != domain.ErrInvalidLongitude {
		t.Errorf("expected ErrInvalidLongitude, got %v", err)
	}
}
