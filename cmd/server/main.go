package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pinbrain/gophkeeper/internal/logger"
	"github.com/pinbrain/gophkeeper/internal/server"
	"github.com/pinbrain/gophkeeper/internal/server/config"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	logger, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	server, err := server.NewServer(cfg, logger)
	if err != nil {
		panic(err)
	}
	logger.WithFields(logrus.Fields{
		"address":  cfg.ServerAddress,
		"logLevel": logger.Level.String(),
	}).Info("Starting server")
	if err = server.Run(); err != nil {
		panic(err)
	}
}
