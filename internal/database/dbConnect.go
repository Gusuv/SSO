package db

import (
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MustDbConnect(log *slog.Logger, dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("Connection failed", slog.Any("error", err))
		os.Exit(1)
	}
	return db

}
