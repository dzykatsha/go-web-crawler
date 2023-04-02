package settings

import "github.com/spf13/viper"

const ASYNQ_CONCURRENCY = "ASYNQ_CONCURRENCY"

type AsynqSettings struct {
	Concurrency int
}

func NewAsynqSettingsFromEnv() *AsynqSettings {
	return &AsynqSettings{
		Concurrency: viper.GetInt(ASYNQ_CONCURRENCY),
	}
}