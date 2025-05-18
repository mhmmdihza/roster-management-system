package main

import (
	"fmt"
	"log"
	"os"
	"payd/handler"
	"payd/util"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
)

func main() {
	initLog()
	logrus := util.Log()
	logrus.Debug("debug mode")

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	if port == ":" {
		port = ":8080"
	}
	logrus.WithField("port", port).Info("starting...")
	httpHandler := handler.NewHandler()
	if err := httpHandler.Run(port); err != nil {
		logrus.Fatal(err)
	}
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
