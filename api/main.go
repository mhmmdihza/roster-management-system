package main

import (
	"fmt"
	"os"
	"payd/handler"
	"payd/services/auth"
	"payd/storage"
	"payd/util"
	"strconv"

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

	st := initStorage()
	authSvc := initAuth(st)

	logrus.WithField("port", port).Info("starting...")
	httpHandler, err := handler.NewHandler(handler.WithAuthSvc(authSvc))
	if err != nil {
		logrus.Fatal(err)
	}
	if err := httpHandler.Run(port); err != nil {
		logrus.Fatal(err)
	}
}

func initAuth(st *storage.Storage) *auth.Auth {
	kratosPubliURL := os.Getenv("KRATO_PUBLIC_URL")
	kratosAdminUrl := os.Getenv("KRATO_ADMIN_URL")

	authSvc, err := auth.NewAuth(st, auth.WithKratosPublicURL(kratosPubliURL), auth.WithKratosAdminURL(kratosAdminUrl))
	if err != nil {
		util.Log().Fatal(err)
	}
	adminEmail, adminPassword, adminName := os.Getenv("BOOTSTRAP_ADMIN_EMAIL"), os.Getenv("BOOTSTRAP_ADMIN_PASSWORD"), os.Getenv("BOOTSTRAP_ADMIN_NAME")
	// skipping bootstraping admin account
	if adminEmail == "" {
		return authSvc
	}
	if adminPassword == "" || adminName == "" {
		util.Log().Fatal("BOOTSTRAP_ADMIN_EMAIL is set, but BOOTSTRAP_ADMIN_PASSWORD or BOOTSTRAP_ADMIN_NAME are missing. All three env are required to bootstrap the admin account.")
	}

	if err := authSvc.BootstrapAdminAccount(adminEmail, adminName, adminPassword); err != nil {
		util.Log().Fatal(err)
	}
	return authSvc
}

func initStorage() *storage.Storage {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSL_MODE")
	migrationDir := os.Getenv("DB_MIGRATION_DIR")
	dbPortInt, _ := strconv.Atoi(dbPort)
	st, err := storage.NewStorage(storage.WithUser(dbUser), storage.WithPassword(dbPassword), storage.WithHost(dbHost),
		storage.WithPort(dbPortInt), storage.WithDbname(dbName), storage.WithSSLMode(sslMode))
	if err != nil {
		util.Log().Fatal(err)
	}
	if err := st.RunMigrations(migrationDir); err != nil {
		util.Log().Fatal(err)
	}
	return st
}

func initLog() {
	logLevelStr := os.Getenv("LOG_LEVEL")
	logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		util.Log().Printf("unknown log level :'%s' , fallback to the info level", logLevelStr)
		logLevel = logrus.InfoLevel
	}
	util.InitLogger(logLevel)
}
