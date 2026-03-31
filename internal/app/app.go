package app

import (
	"log/slog"
	appgrpc "main/internal/app/grpc"
	"main/internal/config"
	"main/internal/repository"
	security "main/internal/security/jwt"
	"main/internal/service"

	"gorm.io/gorm"
)

type App struct {
	GRPCServer *appgrpc.App
}

func New(log *slog.Logger, cfg *config.Config, db *gorm.DB) *App {

	authJWT := security.NewToken(cfg.JWTSecret, cfg.TokenTTL)
	authRepo := repository.NewRepo(db)
	authService := service.New(log, authRepo, cfg.TokenTTL, authJWT)
	grpcServer := appgrpc.New(log, cfg.Grpc.Port, authService)

	return &App{
		GRPCServer: grpcServer,
	}
}
