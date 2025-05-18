package util

import (
	"fmt"
	"path/filepath"
	"runtime"
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
		logger.SetReportCaller(true)
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				return "", fmt.Sprintf("%s:%d", filepath.Base(f.File), f.Line)
			},
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
