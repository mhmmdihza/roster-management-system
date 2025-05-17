package main

import (
	"log"
	"os"
	"payd/util"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
)

func main() {
	initLog()
	logrus := util.Log()
	logrus.Debug("debug mode")
	logrus.Info("starting...")
}
func initLog() {
	logLevelStr := os.Getenv("LOG_LEVEL")
	logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Printf("unknown log level :'%s' , fallback to the info level", logLevelStr)
		logLevel = logrus.InfoLevel
	}
	util.InitLogger(logLevel)
}
