// Package logger реализует настроенный логгер.
package logger

import (
	"github.com/sirupsen/logrus"
)

// NewLogger создает и возвращает новый логгер.
func NewLogger(level string) (*logrus.Logger, error) {
	var err error
	logLvl := logrus.InfoLevel
	if level != "" {
		logLvl, err = logrus.ParseLevel(level)
		if err != nil {
			return nil, err
		}
	}
	logger := logrus.New()
	logger.SetLevel(logLvl)
	logger.SetFormatter(&logrus.JSONFormatter{})
	return logger, nil
}
