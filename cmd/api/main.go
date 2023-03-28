package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dzykatsha/go-web-crawler/internal/api"
	"github.com/dzykatsha/go-web-crawler/internal/settings"
	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)

	redisSettings := settings.NewRedisSettingsFromEnv()
	redisConnection := asynq.RedisClientOpt{
		Addr:     redisSettings.Address(),
		Password: redisSettings.Password,
	}

	client := asynq.NewClient(redisConnection)
	defer client.Close()

	apiSettings := settings.NewAPISettingsFromEnv()

	http.Handle("/load/", api.NewPostLoadURLHandler(client))

	fmt.Printf("Redis: %s\n", redisSettings.Address())
	fmt.Printf("Running on %s\n", apiSettings.Address())
	if err := http.ListenAndServe(apiSettings.Address(), nil); err != nil {
		log.Fatal(err)
	}
}
