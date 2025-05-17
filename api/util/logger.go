package util

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	once   sync.Once
	logger *logrus.Logger
)

// InitLogger initializes the logger with the desired level.
func InitLogger(level logrus.Level) {
	once.Do(func() {
		logger = logrus.New()
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
		logger.SetLevel(level)
	})
}

// It should be called only after InitLogger.
func Log() *logrus.Logger {
	if logger == nil {
		InitLogger(logrus.InfoLevel)
	}
	return logger
}
