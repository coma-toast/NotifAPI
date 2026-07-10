package utils

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// config is the configuration struct
type Config struct {
	// Name of the server running the app
	Name           string
	LogFilePath    string
	PostgresHost   string
	PostgresPort   string
	PostgresUser   string
	PostgresDB     string
	PostgresPass   string
	InstanceID     string
	IPInfoToken    string
	SecretKey      string
	Port           string
	DiscordWebhook string
	JWTKey         string
	DevMode        bool
}

func GetConf(path string) *Config {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Read environment variables with NOTIF_API_ prefix
	viper.SetEnvPrefix("NOTIF_API")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Defaults
	viper.SetDefault("Name", "")
	viper.SetDefault("LogFilePath", "./logs/")
	viper.SetDefault("PostgresHost", "localhost")
	viper.SetDefault("PostgresPort", "5432")
	viper.SetDefault("PostgresUser", "postgres")
	viper.SetDefault("PostgresDB", "postgres")
	viper.SetDefault("PostgresPass", "")
	viper.SetDefault("InstanceID", "")
	viper.SetDefault("IPInfoToken", "")
	viper.SetDefault("SecretKey", "")
	viper.SetDefault("DiscordWebhook", "")
	viper.SetDefault("Port", "10887")
	viper.SetDefault("JWTKey", "")
	viper.SetDefault("DevMode", false)

	// Read config file (optional - env vars take precedence)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using environment variables and defaults")
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
		fmt.Println("InstanceID is required. Set NOTIF_API_INSTANCEID environment variable.")
		os.Exit(1)
	}

	return conf
}