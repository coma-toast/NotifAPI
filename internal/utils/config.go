package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

// config is the configuration struct
type Config struct {
	// Name of the server running the app
	Name              string
	LogFilePath       string
	DBFilePath        string
	InstanceID        string
	SecretKey         string
	Port              string
	DiscordWebhookURL string
	DevMode           bool
}

func GetConf(path string) *Config {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetDefault("Name", "")
	viper.SetDefault("LogFilePath", "./logs/")
	viper.SetDefault("DBFilePath", "./db")
	viper.SetDefault("InstanceID", "")
	viper.SetDefault("SecretKey", "")
	viper.SetDefault("Port", "")
	viper.SetDefault("DiscordWebhookURL", "")
	viper.SetDefault("DevMode", false)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.Unmarshal(&Config{})
			err = viper.SafeWriteConfig()
			if err != nil {
				log.Fatalln("Error writing config", err)
			}

		} else {
			log.Fatalf("Error reading config: %v", err)
		}
	}

	conf := &Config{}
	err := viper.Unmarshal(conf)

	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(conf.LogFilePath); os.IsNotExist(err) {
		os.Mkdir(conf.LogFilePath, 0777)
	}

	if conf.InstanceID == "" {
		fmt.Println("Please fill out configuration values in config.yaml")
		os.Exit(0)
	}

	return conf
}
