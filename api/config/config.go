package config

import (
	"fmt"
	"github.com/mcavoyk/quirk/models"
	"github.com/spf13/viper"
	"log"
)

const Path = "config.toml"

func InitConfig() (*viper.Viper, error) {
	config := viper.New()

	config.SetTypeByDefaultValue(true)
	log.Printf("Reading configuration in from [%s]\n", Path)
	config.SetConfigFile(Path)
	config.WatchConfig()
	err := config.ReadInConfig()
	return config, err
}

func InitDB(config *viper.Viper) (*models.DB, error) {
	dbConnection := fmt.Sprintf("%s:%s@tcp(%s)/quirkdb",
		config.GetString("database.username"),
		config.GetString("database.password"),
		config.GetString("database.address"))

	log.Printf("Attempting to connect to database [%s]\n", dbConnection)

	db, err := models.InitDB(dbConnection + "?charset=utf8&parseTime=True")
	return db, err
}
