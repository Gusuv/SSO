package app

import (
	"log/slog"
	appgrpc "main/internal/app/grpc"
	"time"
)

type App struct {
	GRPCServer *appgrpc.App
}

func New(log *slog.Logger, grpcPort int, storePath string, tokenTTL time.Duration) *App {
	grpc := appgrpc.New(log, grpcPort)

	return &App{
		GRPCServer: grpc,
	}
}
