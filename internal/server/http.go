package server

import (
	"fmt"
	"go-prictic/internal/config"
	"go-prictic/internal/httpresponse"
	"go-prictic/internal/user"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

func Start(cfg *config.Config, db *gorm.DB) error {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.Validator = &CustomValidator{validator: validator.New()}

	if err := db.AutoMigrate(&user.User{}); err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}

	fmt.Println("Database Migrate successfully")

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, httpresponse.Success{
			Code:    http.StatusOK,
			Message: "Welcome to the server!",
			Data:    nil,
		})
	})

	user.Register(e, db, cfg)

	return e.Start(":" + cfg.Port)
}
