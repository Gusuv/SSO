package main

import (
	"main/internal/app"
	"main/internal/config"
	"main/internal/database"
	"main/logger"
)

func main() {

	cfg := config.MustLoad()

	log := logger.AddLogger(cfg.Env)

	db := database.MustDbConnect(log, cfg)

	app := app.New(log, cfg, db)

	go func() {
		if err := app.GRPCServer.Run(); err != nil {
			log.Error("grpc server stopped", "error", err)
		}
	}()

	app.GRPCServer.ReadStop()

	// In development

}
