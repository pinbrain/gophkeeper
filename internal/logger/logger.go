package logger

import (
	"github.com/sirupsen/logrus"
)

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
	// if logLvl == logrus.DebugLevel {
	// 	logger.SetReportCaller(true)
	// }
	return logger, nil
}
