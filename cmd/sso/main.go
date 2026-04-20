package main

import (
	"github.com/Gusuv/sso/internal/app"
	"github.com/Gusuv/sso/internal/config"
	"github.com/Gusuv/sso/internal/database"
	"github.com/Gusuv/sso/logger"
)

func main() {

	cfgPath := config.FetchConfigPath()

	cfg := config.MustLoad(cfgPath)

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
