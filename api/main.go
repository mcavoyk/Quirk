package main

import (
	"context"
	"fmt"
	"log"
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
	logrus.SetOutput(os.Stdout)
	configPath := os.Getenv("CONFIG")
	if configPath == "" {
		configPath = DefaultConfig
	}
	config, err := configuration.InitConfig(configPath)
	if err != nil {
		log.Fatalf("Unable to read configuration: %s", err.Error())
	}

	db, err := configuration.InitDB(config)
	if err != nil {
		logrus.Fatalf("Unable to connect to database: %s", err.Error())
	}
	defer db.Close()

	url := fmt.Sprintf("%s:%d", config.GetString("server.address"), config.GetInt("server.port"))
	log.Printf("Starting server on [%s]", url)
	handler := server.NewRouter(&server.Env{DB: db, Debug: config.GetBool("server.debug_mode")})
	srv := &http.Server{
		Addr:         url,
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
