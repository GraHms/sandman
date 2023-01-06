package config

import (
	"github.com/spf13/viper"
	"log"
)

var config *viper.Viper

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns they configuration struct.
func Init() {
	var err error
	config = viper.New()
	config.SetConfigType("env")
	config.AddConfigPath("../etc/")
	config.AddConfigPath("./etc/")
	config.AddConfigPath("./../../../../../../../etc/")
	config.SetConfigName("config")
	config.AutomaticEnv()

	err = config.ReadInConfig()
	if err != nil {
		log.Fatal("error on parsing configuration file")
	}
}

func GetConfig() *viper.Viper {
	return config
}
