package settings

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	REDIS_HOST     = "REDIS_HOST"
	REDIS_PORT     = "REDIS_PORT"
	REDIS_PASSWORD = "REDIS_PASSWORD"
)

type RedisSettings struct {
	Host     string
	Port     int
	Password string
}

func NewRedisSettingsFromEnv() *RedisSettings {
	return &RedisSettings{
		Host:     viper.GetString(REDIS_HOST),
		Port:     viper.GetInt(REDIS_PORT),
		Password: viper.GetString(REDIS_PASSWORD),
	}
}

func (settings *RedisSettings) Address() string {
	return fmt.Sprintf("%s:%d", settings.Host, settings.Port)
}
