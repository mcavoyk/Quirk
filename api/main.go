package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sirupsen/logrus"

	configuration "github.com/mcavoyk/quirk/api/config"

	"github.com/mcavoyk/quirk/api/server"
)

const DefaultConfig = "config.toml"

func main() {
	log := logrus.New()
	log.SetOutput(os.Stdout)

	configPath := os.Getenv("CONFIG")
	if configPath == "" {
		configPath = DefaultConfig
	}
	config, err := configuration.InitConfig(configPath)
	if err != nil {
		log.Fatalf("Unable to read configuration: %s", err.Error())
	}

	levelStr := config.GetString("server.log_level")
	logLevel, err := logrus.ParseLevel(levelStr)
	if err != nil {
		logrus.Warnf("Unable to parse configuration log_level: %s", levelStr)
		logLevel = logrus.DebugLevel
	}
	log.SetLevel(logLevel)

	db, err := configuration.InitDB(config)
	if err != nil {
		logrus.Fatalf("Unable to setup database: %s", err.Error())
	}
	defer db.Close()

	port := config.GetString("server.port")
	log.Infof("Starting server on port %s", port)
	handler := server.NewRouter(&server.Env{DB: db, Log: log})
	srv := &http.Server{
		Addr:         port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	// Start the server
	go func() {
		log.Fatalf("Shutting down server - %s", srv.ListenAndServe().Error())
	}()

	// Wait for an interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Attempt a graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
