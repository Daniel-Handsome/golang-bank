package utils

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var App *Config

type Config struct {
	DB `mapstructure:",squash"`
}

type DB struct {
	Connection string `mapstructure:"DB_CONNECTION"`
	Host string `mapstructure:"DB_HOST"`
	Port int32 `mapstructure:"DB_PORT"`
	Database string `mapstructure:"DB_DATABASE"`
	Username string `mapstructure:"DB_USERNAME"`
	Password string `mapstructure:"DB_PASSWORD"`
	Token_symmetric_key string `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	Access_token_duration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config,err error) {
	viper.SetConfigFile(path)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
    if err!= nil {
        return
    }

	App = &config

	//查看文件變化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config changed")
	})

	return
}