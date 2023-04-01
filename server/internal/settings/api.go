package settings

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	API_PORT = "API_PORT"
)

type APISettings struct {
	Port int
}

func NewAPISettingsFromEnv() *APISettings {
	return &APISettings{
		Port: viper.GetInt(API_PORT),
	}
}

func (settings *APISettings) Address() string {
	return fmt.Sprintf(
		":%d",
		settings.Port,
	)
}
