package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectionDB(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(
		postgres.New(postgres.Config{
			DSN:                  cfg.DatabaseUrl,
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{
			TranslateError: true,
		},
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}
