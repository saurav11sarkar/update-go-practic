package user

import (
	"errors"
	"go-prictic/internal/httpresponse"
	"go-prictic/internal/user/dto"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateUser(c *echo.Context) error {
	var req dto.UserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Detail:  err.Error(),
		})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Detail:  err.Error(),
		})
	}
	res, err := h.service.CreateUser(req)
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			return c.JSON(http.StatusConflict, httpresponse.Error{
				Code: http.StatusConflict, Message: "User already exists",
			})
		}
		log.Printf("create user failed: %v", err)
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user",
		})
	}
	return c.JSON(http.StatusCreated, httpresponse.Success{
		Code:    http.StatusCreated,
		Message: "User created successfully",
		Data:    res,
	})
}

func (h *Handler) LoginUser(c *echo.Context) error {
	var req dto.UserLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Detail:  err.Error(),
		})
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Validation error",
			Detail:  err.Error(),
		})
	}
	res, err := h.service.LoginUser(req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) || errors.Is(err, ErrUserNotFound) {
			return c.JSON(http.StatusUnauthorized, httpresponse.Error{
				Code: http.StatusUnauthorized, Message: "Invalid email or password",
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to login user",
		})
	}
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		Path:     "/",
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   false, // Set to true when the API is served over HTTPS.
		SameSite: http.SameSiteLaxMode,
	})
	return c.JSON(http.StatusOK, httpresponse.Success{
		Code:    http.StatusOK,
		Message: "User logged in successfully",
		Data:    res,
	})
}

func (h *Handler) LogoutUser(c *echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   false, // Set to true when the API is served over HTTPS.
		SameSite: http.SameSiteLaxMode,
	})
	return c.JSON(http.StatusOK, httpresponse.Success{
		Code:    http.StatusOK,
		Message: "User logged out successfully",
		Data:    nil,
	})
}

func (h *Handler) MyProfile(c *echo.Context) error {
	id, ok := c.Get("id").(string)
	if !ok || id == "" {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code:    http.StatusUnauthorized,
			Message: "Failed to get user id",
		})
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.Error{
				Code: http.StatusNotFound, Message: "User not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get user",
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Code:    http.StatusOK,
		Message: "My profile",
		Data:    user,
	})
}
