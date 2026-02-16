package main

import (
	"log/slog"
	"main/internal/app"
	"main/internal/config"
	"main/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	cfg := config.MustLoad()

	log := logger.AddLogger(cfg.Env)

	application := app.New(log, cfg.Grpc.Port, cfg.DbPath, cfg.TokenTTL)

	go application.GRPCServer.Run()

	stopped := make(chan os.Signal, 1)
	signal.Notify(stopped, syscall.SIGTERM, syscall.SIGINT)

	stopName := <-stopped

	log.Info("Stopping", slog.String("signal", stopName.String()))

	application.GRPCServer.GracefulStop()

	log.Info("Stopped")

	// In development

}
