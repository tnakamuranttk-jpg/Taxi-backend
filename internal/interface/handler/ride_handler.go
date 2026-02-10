// Package handler はHTTPリクエストのハンドリングを担当します
package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tnakamura/taxi-backend/internal/domain"
	"github.com/tnakamura/taxi-backend/internal/interface/handler/dto"
	"github.com/tnakamura/taxi-backend/internal/usecase"
)

// RideHandler は配車関連のHTTPハンドラーです
type RideHandler struct {
	rideUsecase usecase.RideUsecase
}

// NewRideHandler は新しい RideHandler を生成します
func NewRideHandler(u usecase.RideUsecase) *RideHandler {
	return &RideHandler{
		rideUsecase: u,
	}
}

// CreateRide は新しい配車リクエストを作成します
// @Summary 配車リクエスト作成
// @Description 新しい配車リクエストを作成します
// @Tags rides
// @Accept json
// @Produce json
// @Param request body dto.CreateRideRequest true "配車リクエスト"
// @Success 201 {object} dto.RideResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /api/v1/rides [post]
func (h *RideHandler) CreateRide(c *gin.Context) {
	var req dto.CreateRideRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	ride, err := h.rideUsecase.CreateRide(
		req.PassengerID,
		req.PickupLatitude,
		req.PickupLongitude,
		req.DropoffLatitude,
		req.DropoffLongitude,
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "Passenger not found",
			})
		case errors.Is(err, domain.ErrActiveRideExists):
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "conflict",
				Message: "Passenger already has an active ride",
			})
		case errors.Is(err, domain.ErrInvalidLatitude):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid latitude",
			})
		case errors.Is(err, domain.ErrInvalidLongitude):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid longitude",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to create ride",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, dto.ToRideResponse(ride))
}

// GetRide は指定されたIDの配車を取得します
// @Summary 配車情報取得
// @Description IDを指定して配車情報を取得します
// @Tags rides
// @Accept json
// @Produce json
// @Param id path string true "配車ID"
// @Success 200 {object} dto.RideResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/rides/{id} [get]
func (h *RideHandler) GetRide(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Ride ID is required",
		})
		return
	}

	ride, err := h.rideUsecase.GetRide(id)
	if err != nil {
		if errors.Is(err, domain.ErrRideNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "Ride not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get ride",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ToRideResponse(ride))
}

// AcceptRide はドライバーが配車を承諾します
// @Summary 配車承諾
// @Description ドライバーが配車リクエストを承諾します
// @Tags rides
// @Accept json
// @Produce json
// @Param id path string true "配車ID"
// @Param request body dto.AcceptRideRequest true "承諾リクエスト"
// @Success 200 {object} dto.RideResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /api/v1/rides/{id}/accept [put]
func (h *RideHandler) AcceptRide(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Ride ID is required",
		})
		return
	}

	var req dto.AcceptRideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	ride, err := h.rideUsecase.AcceptRide(id, req.DriverID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRideNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "Ride not found",
			})
		case errors.Is(err, domain.ErrDriverNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "Driver not found",
			})
		case errors.Is(err, domain.ErrDriverNotAvailable):
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "conflict",
				Message: "Driver is not available",
			})
		case errors.Is(err, domain.ErrActiveRideExists):
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "conflict",
				Message: "Driver already has an active ride",
			})
		case errors.Is(err, domain.ErrInvalidRideStatus):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "bad_request",
				Message: "Ride cannot be accepted in current status",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to accept ride",
			})
		}
		return
	}

	c.JSON(http.StatusOK, dto.ToRideResponse(ride))
}

// ArriveAtPickup はドライバーが乗車地点に到着したことを記録します
// @Summary 到着通知
// @Description ドライバーが乗車地点に到着したことを記録します
// @Tags rides
// @Accept json
// @Produce json
// @Param id path string true "配車ID"
// @Success 200 {object} dto.RideResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/rides/{id}/arrive [put]
func (h *RideHandler) ArriveAtPickup(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Ride ID is required",
		})
		return
	}

	ride, err := h.rideUsecase.ArriveAtPickup(id)
	if err != nil {
		h.handleRideStatusError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToRideResponse(ride))
}

// StartRide は乗車を開始します
// @Summary 乗車開始
// @Description 乗車を開始します
// @Tags rides
// @Accept json
// @Produce json
// @Param id path string true "配車ID"
// @Success 200 {object} dto.RideResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/rides/{id}/start [put]
func (h *RideHandler) StartRide(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Ride ID is required",
		})
		return
	}

	ride, err := h.rideUsecase.StartRide(id)
	if err != nil {
		h.handleRideStatusError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToRideResponse(ride))
}

// CompleteRide は乗車を完了します
// @Summary 乗車完了
// @Description 乗車を完了します
// @Tags rides
// @Accept json
// @Produce json
// @Param id path string true "配車ID"
// @Success 200 {object} dto.RideResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/rides/{id}/complete [put]
func (h *RideHandler) CompleteRide(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Ride ID is required",
		})
		return
	}

	ride, err := h.rideUsecase.CompleteRide(id)
	if err != nil {
		h.handleRideStatusError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToRideResponse(ride))
}

// CancelRide は配車をキャンセルします
// @Summary 配車キャンセル
// @Description 配車をキャンセルします
// @Tags rides
// @Accept json
// @Produce json
// @Param id path string true "配車ID"
// @Success 200 {object} dto.RideResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/rides/{id}/cancel [put]
func (h *RideHandler) CancelRide(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Ride ID is required",
		})
		return
	}

	ride, err := h.rideUsecase.CancelRide(id)
	if err != nil {
		h.handleRideStatusError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToRideResponse(ride))
}

// GetPassengerRides は乗客の配車履歴を取得します
// @Summary 乗客の配車履歴取得
// @Description 指定した乗客の配車履歴を取得します
// @Tags rides
// @Accept json
// @Produce json
// @Param passenger_id path string true "乗客ID"
// @Success 200 {array} dto.RideResponse
// @Router /api/v1/users/{passenger_id}/rides [get]
func (h *RideHandler) GetPassengerRides(c *gin.Context) {
	passengerID := c.Param("passenger_id")
	if passengerID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Passenger ID is required",
		})
		return
	}

	rides, err := h.rideUsecase.GetPassengerRides(passengerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get rides",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ToRideResponseList(rides))
}

// GetDriverRides はドライバーの配車履歴を取得します
// @Summary ドライバーの配車履歴取得
// @Description 指定したドライバーの配車履歴を取得します
// @Tags rides
// @Accept json
// @Produce json
// @Param driver_id path string true "ドライバーID"
// @Success 200 {array} dto.RideResponse
// @Router /api/v1/drivers/{driver_id}/rides [get]
func (h *RideHandler) GetDriverRides(c *gin.Context) {
	driverID := c.Param("driver_id")
	if driverID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Driver ID is required",
		})
		return
	}

	rides, err := h.rideUsecase.GetDriverRides(driverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get rides",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ToRideResponseList(rides))
}

// handleRideStatusError はステータス関連のエラーを共通処理します
func (h *RideHandler) handleRideStatusError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrRideNotFound):
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "Ride not found",
		})
	case errors.Is(err, domain.ErrInvalidRideStatus):
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Invalid ride status transition",
		})
	default:
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to update ride",
		})
	}
}

// EstimateFare は乗降車地点から料金を見積もります
// @Summary 料金見積もり
// @Description 乗車地点と降車地点を指定して料金を見積もります
// @Tags rides
// @Accept json
// @Produce json
// @Param request body dto.EstimateRideRequest true "料金見積もりリクエスト"
// @Success 200 {object} dto.FareEstimateResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/rides/estimate [post]
func (h *RideHandler) EstimateFare(c *gin.Context) {
	var req dto.EstimateRideRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	fare, err := h.rideUsecase.EstimateFare(
		req.PickupLatitude,
		req.PickupLongitude,
		req.DropoffLatitude,
		req.DropoffLongitude,
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidLatitude):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid latitude",
			})
		case errors.Is(err, domain.ErrInvalidLongitude):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid longitude",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to estimate fare",
			})
		}
		return
	}

	c.JSON(http.StatusOK, dto.ToFareEstimateResponse(fare))
}
