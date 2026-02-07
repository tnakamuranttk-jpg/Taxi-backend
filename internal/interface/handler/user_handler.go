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

// UserHandler はユーザー関連のHTTPハンドラーです
type UserHandler struct {
	userUsecase usecase.UserUsecase
}

// NewUserHandler は新しい UserHandler を生成します
func NewUserHandler(u usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: u,
	}
}

// GetUser は指定されたIDのユーザーを取得します
// @Summary ユーザー取得
// @Description IDを指定してユーザー情報を取得します
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ユーザーID"
// @Success 200 {object} dto.UserResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "User ID is required",
		})
		return
	}

	user, err := h.userUsecase.GetUser(id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "not_found",
				Message: "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get user",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ToUserResponse(user))
}

// CreateUser は新しいユーザーを作成します
// @Summary ユーザー作成
// @Description 新しいユーザーを登録します
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "ユーザー作成リクエスト"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	user, err := h.userUsecase.CreateUser(req.Name, req.Email)
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
		case errors.Is(err, domain.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "conflict",
				Message: "User with this email already exists",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to create user",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, dto.ToUserResponse(user))
}

// GetAllUsers は全てのユーザーを取得します
// @Summary ユーザー一覧取得
// @Description 登録されている全ユーザーを取得します
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} dto.UserResponse
// @Router /api/v1/users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userUsecase.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get users",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ToUserResponseList(users))
}
