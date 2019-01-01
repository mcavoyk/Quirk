package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"./models"
	"./server"
	"github.com/spf13/viper"
)

const Config = "config.toml"

func main() {
	log.Printf("Reading configuration in from [%s]\n", Config)
	config, err := initConfig()
	if err != nil {
		log.Fatalf("Unable to read configuration: %s", err.Error())
	}

	var db *models.DB

	for true {
		dbConnection := fmt.Sprintf("%s:%s@tcp(%s)/quirkdb",
			config.GetString("database.username"),
			config.GetString("database.password"),
			config.GetString("database.address"))

		log.Printf("Attempting to connect to database [%s]\n", dbConnection)

		db, err = models.InitDB(dbConnection + "?charset=utf8&parseTime=True")
		if err == nil {
			break
		}
		log.Printf("Unable to connect to database: %s", err.Error())
		time.Sleep(5 * time.Second)
	}
	defer db.Close()

	url := fmt.Sprintf("%s:%d", config.GetString("server.address"), config.GetInt("server.port"))
	log.Printf("Starting server on [%s]", url)
	handler := server.NewRouter(db)
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

func initConfig() (*viper.Viper, error) {
	config := viper.New()

	config.SetTypeByDefaultValue(true)
	config.SetConfigFile(Config)
	config.WatchConfig()
	err := config.ReadInConfig()
	return config, err
}
