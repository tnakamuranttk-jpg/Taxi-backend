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

// DriverHandler はドライバー関連のHTTPハンドラーです
type DriverHandler struct {
	driverUsecase usecase.DriverUsecase
}

// NewDriverHandler は新しい DriverHandler を生成します
func NewDriverHandler(u usecase.DriverUsecase) *DriverHandler {
	return &DriverHandler{
		driverUsecase: u,
	}
}

// GetDriver は指定されたIDのドライバーを取得します
// @Summary ドライバー取得
// @Description IDを指定してドライバー情報を取得します
// @Tags drivers
// @Accept json
// @Produce json
// @Param id path string true "ドライバーID"
// @Success 200 {object} dto.DriverResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/drivers/{id} [get]
func (h *DriverHandler) GetDriver(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Driver ID is required",
		})
		return
	}

	driver, err := h.driverUsecase.GetDriver(id)
	if err != nil {
		if errors.Is(err, domain.ErrDriverNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "Driver not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get driver",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ToDriverResponse(driver))
}

// CreateDriver は新しいドライバーを作成します
// @Summary ドライバー作成
// @Description 新しいドライバーを登録します
// @Tags drivers
// @Accept json
// @Produce json
// @Param request body dto.CreateDriverRequest true "ドライバー作成リクエスト"
// @Success 201 {object} dto.DriverResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /api/v1/drivers [post]
func (h *DriverHandler) CreateDriver(c *gin.Context) {
	var req dto.CreateDriverRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	driver, err := h.driverUsecase.CreateDriver(req.Name, req.Email, req.LicenseNumber)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmptyName):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Name cannot be empty",
			})
		case errors.Is(err, domain.ErrInvalidEmail):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid email format",
			})
		case errors.Is(err, domain.ErrEmptyLicenseNumber):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "License number cannot be empty",
			})
		case errors.Is(err, domain.ErrDriverAlreadyExists):
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "conflict",
				Message: "Driver with this email or license number already exists",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to create driver",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, dto.ToDriverResponse(driver))
}

// GetAllDrivers は全てのドライバーを取得します
// @Summary ドライバー一覧取得
// @Description 登録されている全ドライバーを取得します
// @Tags drivers
// @Accept json
// @Produce json
// @Success 200 {array} dto.DriverResponse
// @Router /api/v1/drivers [get]
func (h *DriverHandler) GetAllDrivers(c *gin.Context) {
	drivers, err := h.driverUsecase.GetAllDrivers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get drivers",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ToDriverResponseList(drivers))
}

// GetAvailableDrivers は配車可能なドライバーを取得します
// @Summary 配車可能ドライバー一覧取得
// @Description 配車可能な（空車状態の）ドライバーを取得します
// @Tags drivers
// @Accept json
// @Produce json
// @Success 200 {array} dto.DriverResponse
// @Router /api/v1/drivers/available [get]
func (h *DriverHandler) GetAvailableDrivers(c *gin.Context) {
	drivers, err := h.driverUsecase.GetAvailableDrivers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get available drivers",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ToDriverResponseList(drivers))
}

// UpdateDriverLocation はドライバーの現在地を更新します
// @Summary ドライバー位置更新
// @Description ドライバーの現在地（緯度・経度）を更新します
// @Tags drivers
// @Accept json
// @Produce json
// @Param id path string true "ドライバーID"
// @Param request body dto.UpdateDriverLocationRequest true "位置更新リクエスト"
// @Success 200 {object} dto.DriverResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/drivers/{id}/location [put]
func (h *DriverHandler) UpdateDriverLocation(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Driver ID is required",
		})
		return
	}

	var req dto.UpdateDriverLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	driver, err := h.driverUsecase.UpdateDriverLocation(id, req.Latitude, req.Longitude)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDriverNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "Driver not found",
			})
		case errors.Is(err, domain.ErrInvalidLatitude):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Latitude must be between -90 and 90",
			})
		case errors.Is(err, domain.ErrInvalidLongitude):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Longitude must be between -180 and 180",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to update driver location",
			})
		}
		return
	}

	c.JSON(http.StatusOK, dto.ToDriverResponse(driver))
}

// UpdateDriverStatus はドライバーのステータスを更新します
// @Summary ドライバーステータス更新
// @Description ドライバーの稼働ステータスを更新します（available/busy/offline）
// @Tags drivers
// @Accept json
// @Produce json
// @Param id path string true "ドライバーID"
// @Param request body dto.UpdateDriverStatusRequest true "ステータス更新リクエスト"
// @Success 200 {object} dto.DriverResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/drivers/{id}/status [put]
func (h *DriverHandler) UpdateDriverStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Driver ID is required",
		})
		return
	}

	var req dto.UpdateDriverStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	status := domain.DriverStatus(req.Status)
	driver, err := h.driverUsecase.UpdateDriverStatus(id, status)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDriverNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "Driver not found",
			})
		case errors.Is(err, domain.ErrInvalidDriverStatus):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid driver status",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to update driver status",
			})
		}
		return
	}

	c.JSON(http.StatusOK, dto.ToDriverResponse(driver))
}
