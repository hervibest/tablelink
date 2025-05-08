package config

import "github.com/spf13/viper"

type Config struct {
	PgURL     string
	RedisAddr string
	PortAuth  string
	PortUsers string
}

func Load() (*Config, error) {
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	cfg := &Config{
		PgURL:     viper.GetString("PG_URL"),
		RedisAddr: viper.GetString("REDIS_ADDR"),
		PortAuth:  viper.GetString("PORT_AUTH"),
		PortUsers: viper.GetString("PORT_USERS"),
	}
	return cfg, nil

}
