package app

import (
	"log/slog"

	appgrpc "github.com/Gusuv/sso/internal/app/grpc"
	"github.com/Gusuv/sso/internal/config"
	"github.com/Gusuv/sso/internal/repository"
	"github.com/Gusuv/sso/internal/security/hash"
	security "github.com/Gusuv/sso/internal/security/jwt"
	"github.com/Gusuv/sso/internal/service"

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
