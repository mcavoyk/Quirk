package config

import (
	"log"

	"github.com/spf13/viper"
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
