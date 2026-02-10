// Package domain はビジネスロジックの中核となるドメイン層を定義します
package domain

import "math"

// 料金定数（日本の一般的なタクシー料金体系を参考）
const (
	// BaseFare は初乗り料金（円）
	BaseFare = 500
	// BaseDistance は初乗り距離（km）
	BaseDistance = 1.0
	// PerKmRate は加算距離あたりの料金（円/km）
	PerKmRate = 300
	// MinimumFare は最低料金（円）
	MinimumFare = 500
	// EarthRadiusKm は地球の平均半径（km）
	EarthRadiusKm = 6371.0
)

// FareBreakdown は料金の内訳を表します
type FareBreakdown struct {
	BaseFare     int     `json:"base_fare"`     // 初乗り料金（円）
	DistanceFare int     `json:"distance_fare"` // 距離加算料金（円）
	TotalFare    int     `json:"total_fare"`    // 合計料金（円）
	DistanceKm   float64 `json:"distance_km"`   // 走行距離（km）
}

// CalculateDistance は2地点間の距離をHaversine公式で計算します（km）
func CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	// 角度をラジアンに変換
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLng := (lng2 - lng1) * math.Pi / 180

	// Haversine公式
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadiusKm * c
}

// CalculateFare は走行距離から料金を計算します
func CalculateFare(distanceKm float64) *FareBreakdown {
	baseFare := BaseFare
	distanceFare := 0

	// 初乗り距離を超えた分に対して加算
	if distanceKm > BaseDistance {
		extraDistance := distanceKm - BaseDistance
		distanceFare = int(math.Ceil(extraDistance * PerKmRate))
	}

	totalFare := baseFare + distanceFare
	if totalFare < MinimumFare {
		totalFare = MinimumFare
	}

	return &FareBreakdown{
		BaseFare:     baseFare,
		DistanceFare: distanceFare,
		TotalFare:    totalFare,
		DistanceKm:   math.Round(distanceKm*100) / 100, // 小数点2桁に丸め
	}
}

// CalculateRideFare は配車リクエストの乗降車地点から料金を計算します
func CalculateRideFare(pickupLat, pickupLng, dropoffLat, dropoffLng float64) *FareBreakdown {
	distance := CalculateDistance(pickupLat, pickupLng, dropoffLat, dropoffLng)
	return CalculateFare(distance)
}
