package main

import (
	"payd/util"

	"github.com/sirupsen/logrus"
)

func main() {
	util.InitLogger(logrus.DebugLevel)

	log := util.Log()
	log.Info("starting...")
}
