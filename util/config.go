package util

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sbbullet/to-do/logger"
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerHost          string        `mapstructure:"SERVER_HOST"`
	ServerPort          string        `mapstructure:"SERVER_PORT"`
	SymmetricKey        string        `mapstructure:"SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(fileName string, fileType string, path string) *Config {
	viper.SetConfigName(fileName)
	viper.SetConfigType(fileType)
	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.Panic(err.Error())
	}

	// Set defaults for the config variables
	config := &Config{
		DBDriver:   "sqlite3",
		DBSource:   "todo.db",
		ServerHost: "0.0.0.0",
		ServerPort: "5000",
	}

	// Unmarshal and override config
	err := viper.Unmarshal(config)
	if err != nil {
		logger.Panic(err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		logger.Panic(fmt.Sprintf("Missing required attributes %v", err))
	}

	return config
}
