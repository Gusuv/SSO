package main

import (
	"main/internal/app"
	"main/internal/config"
	"main/logger"
)

func main() {

	cfg := config.MustLoad()

	log := logger.AddLogger(cfg.Env)

	app := app.New(log, cfg.Grpc.Port, cfg.DbPath, cfg.TokenTTL)

	go app.GRPCServer.Run()

	app.GRPCServer.ReadStop()

	// In development

}
