// Package domain はビジネスロジックの中核となるドメイン層を定義します
package domain

import (
	"math"
	"testing"
)

func TestCalculateDistance_SamePoint(t *testing.T) {
	// 同一地点の距離は0
	dist := CalculateDistance(35.6762, 139.6503, 35.6762, 139.6503)
	if dist != 0 {
		t.Errorf("expected 0, got %f", dist)
	}
}

func TestCalculateDistance_TokyoStationToShibuya(t *testing.T) {
	// 東京駅 → 渋谷駅（約5.4km）
	dist := CalculateDistance(35.6812, 139.7671, 35.6580, 139.7016)
	// 許容誤差 ±1km（直線距離なので実際の道路距離とは異なる）
	if dist < 4.0 || dist > 7.0 {
		t.Errorf("expected distance between 4.0 and 7.0 km, got %f", dist)
	}
}

func TestCalculateDistance_TokyoToOsaka(t *testing.T) {
	// 東京 → 大阪（約400km）
	dist := CalculateDistance(35.6762, 139.6503, 34.6937, 135.5023)
	if dist < 380 || dist > 420 {
		t.Errorf("expected distance between 380 and 420 km, got %f", dist)
	}
}

func TestCalculateFare_ShortDistance(t *testing.T) {
	// 初乗り距離以内（0.5km）
	fare := CalculateFare(0.5)

	if fare.BaseFare != BaseFare {
		t.Errorf("expected base fare %d, got %d", BaseFare, fare.BaseFare)
	}
	if fare.DistanceFare != 0 {
		t.Errorf("expected distance fare 0, got %d", fare.DistanceFare)
	}
	if fare.TotalFare != MinimumFare {
		t.Errorf("expected total fare %d, got %d", MinimumFare, fare.TotalFare)
	}
}

func TestCalculateFare_ExactBaseDistance(t *testing.T) {
	// ちょうど初乗り距離（1.0km）
	fare := CalculateFare(BaseDistance)

	if fare.DistanceFare != 0 {
		t.Errorf("expected distance fare 0 for base distance, got %d", fare.DistanceFare)
	}
	if fare.TotalFare != BaseFare {
		t.Errorf("expected total fare %d, got %d", BaseFare, fare.TotalFare)
	}
}

func TestCalculateFare_MediumDistance(t *testing.T) {
	// 5km走行: 初乗り500 + (5-1)*300 = 500 + 1200 = 1700
	fare := CalculateFare(5.0)

	expectedDistanceFare := int(math.Ceil(4.0 * PerKmRate))
	expectedTotal := BaseFare + expectedDistanceFare

	if fare.DistanceFare != expectedDistanceFare {
		t.Errorf("expected distance fare %d, got %d", expectedDistanceFare, fare.DistanceFare)
	}
	if fare.TotalFare != expectedTotal {
		t.Errorf("expected total fare %d, got %d", expectedTotal, fare.TotalFare)
	}
}

func TestCalculateFare_LongDistance(t *testing.T) {
	// 20km走行: 初乗り500 + (20-1)*300 = 500 + 5700 = 6200
	fare := CalculateFare(20.0)

	expectedDistanceFare := int(math.Ceil(19.0 * PerKmRate))
	expectedTotal := BaseFare + expectedDistanceFare

	if fare.TotalFare != expectedTotal {
		t.Errorf("expected total fare %d, got %d", expectedTotal, fare.TotalFare)
	}
}

func TestCalculateFare_DistanceKmRounding(t *testing.T) {
	// 小数点2桁に丸められることを確認
	fare := CalculateFare(3.456789)

	if fare.DistanceKm != 3.46 {
		t.Errorf("expected distance_km 3.46, got %f", fare.DistanceKm)
	}
}

func TestCalculateRideFare_Integration(t *testing.T) {
	// 東京駅 → 渋谷駅
	fare := CalculateRideFare(35.6812, 139.7671, 35.6580, 139.7016)

	if fare.TotalFare <= 0 {
		t.Error("expected positive fare")
	}
	if fare.DistanceKm <= 0 {
		t.Error("expected positive distance")
	}
	if fare.BaseFare != BaseFare {
		t.Errorf("expected base fare %d, got %d", BaseFare, fare.BaseFare)
	}
	// 東京→渋谷は約5km、料金は1500〜2500円程度が想定される
	if fare.TotalFare < 1000 || fare.TotalFare > 3000 {
		t.Errorf("fare seems unreasonable for Tokyo to Shibuya: %d yen", fare.TotalFare)
	}
}

func TestCalculateFare_ZeroDistance(t *testing.T) {
	fare := CalculateFare(0)

	if fare.TotalFare != MinimumFare {
		t.Errorf("expected minimum fare %d for zero distance, got %d", MinimumFare, fare.TotalFare)
	}
}
