package main

import (
	"main/internal/app"
	"main/internal/config"
	db "main/internal/database"
	"main/logger"
)

func main() {

	cfg := config.MustLoad()

	log := logger.AddLogger(cfg.Env)

	db := db.MustDbConnect(log, cfg.DbDsn)

	app := app.New(log, cfg.Grpc.Port, cfg.TokenTTL, db)

	go app.GRPCServer.Run()

	app.GRPCServer.ReadStop()

	// In development

}
