package app

import (
	"log/slog"
	appgrpc "main/internal/app/grpc"
	"main/internal/config"
	"main/internal/repository"
	"main/internal/security/hash"
	security "main/internal/security/jwt"
	"main/internal/service"

	"gorm.io/gorm"
)

type App struct {
	GRPCServer *appgrpc.App
}

func New(log *slog.Logger, cfg *config.Config, db *gorm.DB) *App {
	authHash := hash.NewHash(cfg.HMACSecret)
	authJWT := security.NewToken(cfg.JWTSecret, cfg.TokenTTL)
	authRepo := repository.NewRepo(db)
	authService := service.New(log, authRepo, authJWT, authHash)
	grpcServer := appgrpc.New(log, cfg.Grpc.Port, authService)

	return &App{
		GRPCServer: grpcServer,
	}
}
