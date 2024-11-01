package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pinbrain/gophkeeper/internal/logger"
	"github.com/pinbrain/gophkeeper/internal/server"
	"github.com/pinbrain/gophkeeper/internal/server/config"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const (
	timeoutServerShutdown = time.Second * 5
	timeoutShutdown       = time.Second * 10
)

func main() {
	rootCtx, cancelCtx := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)
	defer cancelCtx()

	g, ctx := errgroup.WithContext(rootCtx)
	// нештатное завершение программы по таймауту
	// происходит, если после завершения контекста
	// приложение не смогло завершиться за отведенный промежуток времени
	context.AfterFunc(ctx, func() {
		afterCtx, afterCancelCtx := context.WithTimeout(context.Background(), timeoutShutdown)
		defer afterCancelCtx()

		<-afterCtx.Done()
		log.Fatal("failed to gracefully shutdown the service")
	})

	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	logger, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	server, err := server.NewServer(ctx, cfg, logger)
	if err != nil {
		panic(err)
	}
	logger.WithFields(logrus.Fields{
		"address":  cfg.ServerAddress,
		"logLevel": logger.Level.String(),
	}).Info("Starting server")

	g.Go(func() (err error) {
		defer func() {
			errRec := recover()
			if errRec != nil {
				err = fmt.Errorf("a panic occurred: %v", errRec)
			}
		}()
		return server.Run()
	})

	// Отслеживаем успешное завершение работы сервера
	g.Go(func() error {
		defer logger.Info("Service has been shutdown")

		<-ctx.Done()
		logger.Info("Gracefully shutting down service...")

		if err = server.Stop(); err != nil {
			logger.Errorf("An error occurred stopping server: %v", err)
		}
		logger.Info("Server fully stopped")

		return nil
	})

	if err = g.Wait(); err != nil {
		panic(err)
	}
}
