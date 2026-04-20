package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/Gusuv/sso/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dbConnectOp = "db.connect"

func MustDbConnect(log *slog.Logger, cfg *config.Config) *gorm.DB {
	dbLog := log.With("op", dbConnectOp)
	dbLog.Info("Connecting to database")

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{TranslateError: true})
	if err != nil {
		dbLog.Error("Connection failed", "error", err)
		panic(fmt.Errorf("Database connect failed: %w", err))
	}
	dbSQL, err := db.DB()
	if err != nil {
		dbLog.Warn("Can`t get SQL.db: %w", "error", err)
	}
	checkConnection(dbSQL, dbLog)

	connectionPool(dbSQL, dbLog, cfg.Env)
	dbLog.Info("Database connected")
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
