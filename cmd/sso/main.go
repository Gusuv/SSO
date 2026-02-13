package main

import (
	"main/internal/app"
	"main/internal/config"
	"main/logger"
)

func main() {

	cfg := config.MustLoad()

	log := logger.AddLogger(cfg.Env)

	application := app.New(log, cfg.Grpc.Port, cfg.DbPath, cfg.TokenTTL)

	application.GRPCServer.Run()

	// In development

}
