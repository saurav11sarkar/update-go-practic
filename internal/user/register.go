package user

import (
	"go-prictic/internal/auth"
	"go-prictic/internal/config"
	"go-prictic/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func Register(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	repository := NewRepository(db)
	accessTokens := auth.NewJwtService(cfg.JwtAccessToken, cfg.JwtTimeDuration)
	refreshTokens := auth.NewJwtService(cfg.JwtRefreshToken, cfg.JwtRefreshTimeDuration)
	service := NewService(repository, accessTokens, refreshTokens)
	handler := NewHandler(service)
	authMiddleware := middlewares.NewAuthMiddleware(accessTokens)

	api := e.Group("/api/v1/auth")
	api.POST("/register", handler.CreateUser)
	api.POST("/login", handler.LoginUser)
	api.POST("/logout", handler.LogoutUser)
	api.GET("/me", handler.MyProfile, authMiddleware.AuthMiddleware("user", "admin"))
}
