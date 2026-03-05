package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"main/internal/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MustDbConnect(log *slog.Logger, cfg config.Config) *gorm.DB {
	log.Info("Connecting to database")
	db, err := gorm.Open(postgres.Open(cfg.DbDsn), &gorm.Config{})
	if err != nil {
		log.Error("Connection failed", "error", err)
		panic(fmt.Errorf("Database connect failed: %w", err))
	}
	dbSQL, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("Can`t get SQL.db: %w", err))
	}

	checkConnection(dbSQL, log)

	connectionPool(dbSQL, log, cfg.Env)
	log.Info("Database connected")
	return db

}

func checkConnection(db *sql.DB, log *slog.Logger) {
	if err := db.Ping(); err != nil {
		log.Error("Database ping failed", "error", err)
		panic(fmt.Errorf("Database ping failed: %w", err))
	}

}

func connectionPool(db *sql.DB, log *slog.Logger, env string) {

	switch env {
	case "local":
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(time.Hour)
		db.SetMaxOpenConns(20)
		db.SetConnMaxIdleTime(5 * time.Minute)

	case "prod":
		db.SetMaxIdleConns(25)
		db.SetConnMaxLifetime(5 * time.Hour)
		db.SetMaxOpenConns(100)
		db.SetConnMaxIdleTime(30 * time.Minute)

	default:
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(time.Hour)
		db.SetMaxOpenConns(20)
		db.SetConnMaxIdleTime(5 * time.Minute)
		log.Warn("env is not defined", "env", env)

	}
}
