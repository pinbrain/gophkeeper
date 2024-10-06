package main

import (
	"github.com/pinbrain/gophkeeper/internal/logger"
	"github.com/pinbrain/gophkeeper/internal/server"
	"github.com/pinbrain/gophkeeper/internal/server/config"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.ServerConfig{
		ServerAddress: ":8080",
		LogLevel:      "debug",
		DSN:           "postgresql://keeper:keeper@192.168.0.27:5412/keeper",
	}

	if err := logger.Initialize(cfg.LogLevel); err != nil {
		panic(err)
	}

	server, err := server.NewServer(cfg)
	if err != nil {
		panic(err)
	}
	logger.Log.WithFields(logrus.Fields{
		"addr":     cfg.ServerAddress,
		"logLevel": cfg.LogLevel,
	}).Info("Starting server")
	if err = server.Run(); err != nil {
		panic(err)
	}
}
