package settings

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	MONGO_HOST       = "MONGO_HOST"
	MONGO_PORT       = "MONGO_PORT"
	MONGO_USERNAME   = "MONGO_USERNAME"
	MONGO_PASSWORD   = "MONGO_PASSWORD"
	MONGO_DATABASE   = "MONGO_DATABASE"
	MONGO_COLLECTION = "MONGO_COLLECTION"
)

type MongoSettings struct {
	Host       string
	Port       int
	Username   string
	Password   string
	Database   string
	Collection string
}

func NewMongoSettingsFromEnv() *MongoSettings {
	return &MongoSettings{
		Host:       viper.GetString(MONGO_HOST),
		Port:       viper.GetInt(MONGO_PORT),
		Username:   viper.GetString(MONGO_USERNAME),
		Password:   viper.GetString(MONGO_PASSWORD),
		Database:   viper.GetString(MONGO_DATABASE),
		Collection: viper.GetString(MONGO_COLLECTION),
	}
}

func (settings *MongoSettings) ConnectionURL() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%d",
		settings.Username,
		settings.Password,
		settings.Host,
		settings.Port,
	)
}
