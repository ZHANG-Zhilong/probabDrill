package conf

import (
	"github.com/spf13/viper"
	"log"
)

func LoadConfig(path string) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	if path == "" {
		viper.AddConfigPath(".")
	} else {
		viper.AddConfigPath(path) // optionally look for config in the working directory
	}
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalln("Config file not found")
		} else {
			log.Fatalln("Config file was found but another error was produced ")
		}
	}
	viper.WatchConfig()

}
