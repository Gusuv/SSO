package app

import (
	"log/slog"
	appgrpc "main/internal/app/grpc"
	"main/internal/repository"
	"main/internal/service"
	"time"

	"gorm.io/gorm"
)

type App struct {
	GRPCServer *appgrpc.App
}

func New(log *slog.Logger, grpcPort int, tokenTTL time.Duration, db *gorm.DB) *App {

	authRepo := repository.NewRepo(db)
	authService := service.New(log, authRepo, tokenTTL)
	grpc := appgrpc.New(log, grpcPort, authService)

	return &App{
		GRPCServer: grpc,
	}
}
