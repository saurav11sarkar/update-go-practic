package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string
	DatabaseUrl            string
	JwtAccessToken         string
	JwtTimeDuration        time.Duration
	JwtRefreshToken        string
	JwtRefreshTimeDuration time.Duration
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	required := []string{"PORT", "DATABASE_URL", "JWT_ACCESS_TOKEN", "JWT_REFRESH_TOKEN"}
	for _, key := range required {
		if os.Getenv(key) == "" {
			return nil, fmt.Errorf("%s is required", key)
		}
	}

	jwtDurationMinutes, err := strconv.Atoi(os.Getenv("JWT_TIME_DURATION"))
	if err != nil || jwtDurationMinutes <= 0 {
		return nil, fmt.Errorf("JWT_TIME_DURATION must be a positive number of minutes")
	}
	jwtRefreshDurationDays, err := strconv.Atoi(os.Getenv("JWT_REFRESH_TIME_DURATION"))
	if err != nil || jwtRefreshDurationDays <= 0 {
		return nil, fmt.Errorf("JWT_REFRESH_TIME_DURATION must be a positive number of days")
	}

	return &Config{
		Port:                   os.Getenv("PORT"),
		DatabaseUrl:            os.Getenv("DATABASE_URL"),
		JwtAccessToken:         os.Getenv("JWT_ACCESS_TOKEN"),
		JwtTimeDuration:        time.Duration(jwtDurationMinutes) * time.Minute,
		JwtRefreshToken:        os.Getenv("JWT_REFRESH_TOKEN"),
		JwtRefreshTimeDuration: time.Duration(jwtRefreshDurationDays) * 24 * time.Hour,
	}, nil
}
