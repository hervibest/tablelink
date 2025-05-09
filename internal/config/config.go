package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	PgURL     string
	RedisAddr string
	PortAuth  string
	PortUsers string
}

func Load(logger *logrus.Logger) (*Config, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	viper.AutomaticEnv()

	cfg := &Config{
		PgURL:     viper.GetString("PG_URL"),
		RedisAddr: viper.GetString("REDIS_ADDR"),
		PortAuth:  viper.GetString("PORT_AUTH"),
		PortUsers: viper.GetString("PORT_USERS"),
	}
	return cfg, nil
}
