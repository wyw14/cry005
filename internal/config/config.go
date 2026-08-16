package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Address        string
	RequestTimeout time.Duration
	DataPath       string
}

func Load() Config {
	viper.SetDefault("address", ":8080")
	viper.SetDefault("request_timeout", "5s")
	viper.SetDefault("data_path", "./data")
	viper.AutomaticEnv()
	return Config{Address: viper.GetString("address"), RequestTimeout: viper.GetDuration("request_timeout"), DataPath: viper.GetString("data_path")}
}
