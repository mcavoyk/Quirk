package config

import (
	"fmt"
	"github.com/mcavoyk/quirk/api/models"
	"github.com/spf13/viper"
	"log"
)

var Path string

func InitConfig(path string) (*viper.Viper, error) {
	Path = path
	config := viper.New()

	config.SetTypeByDefaultValue(true)
	log.Printf("Reading configuration in from [%s]\n", Path)
	config.SetConfigFile(Path)
	err := config.ReadInConfig()
	if err != nil {
		return nil, err
	}
	config.AutomaticEnv()
	return config, nil
}

func InitDB(config *viper.Viper) (*models.DB, error) {
	dbConnection := fmt.Sprintf("%s:%s@tcp(%s)/quirkdb",
		config.GetString("database.username"),
		config.GetString("database.password"),
		config.GetString("database.address"))

	log.Printf("Attempting to connect to database [%s]\n", dbConnection)

	db, err := models.InitDB(dbConnection + "?parseTime=True&charset=utf8mb4&collation=utf8mb4_unicode_ci")
	return db, err
}
