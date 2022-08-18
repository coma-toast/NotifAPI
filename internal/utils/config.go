package utils

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// config is the configuration struct
type Config struct {
	LogFilePath string
	DBFilePath  string
	InstanceID  string
	SecretKey   string
	Port        string
	DevMode     bool
}

func GetConf() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetDefault("LogFilePath", "./logs/")
	viper.SetDefault("InstanceID", "")
	viper.SetDefault("SecretKey", "")
	viper.SetDefault("Port", "")
	viper.SetDefault("DevMode", false)

	err := viper.ReadInConfig()

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
	err = viper.Unmarshal(conf)

	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(conf.LogFilePath); os.IsNotExist(err) {
		os.Mkdir(conf.LogFilePath, 0777)
	}

	return conf
}