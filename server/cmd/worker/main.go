package main

import (
	"context"
	"time"

	"github.com/dzykatsha/go-web-crawler/internal/crawler/load"
	"github.com/dzykatsha/go-web-crawler/internal/settings"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)

	redisSettings := settings.NewRedisSettingsFromEnv()
	mongoSettings := settings.NewMongoSettingsFromEnv()
	asynqSettings := settings.NewAsynqSettingsFromEnv()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoSettings.ConnectionURL()))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to mongodb")
	}

	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	redisConnection := asynq.RedisClientOpt{
		Addr:     redisSettings.Address(),
		Password: redisSettings.Password,
	}

	asynqClient := asynq.NewClient(redisConnection)
	defer asynqClient.Close()

	worker := asynq.NewServer(redisConnection, asynq.Config{
		Concurrency: asynqSettings.Concurrency,
		Queues: map[string]int{
			"critical": 6, // processed 60% of the time
			"default":  3, // processed 30% of the time
			"low":      1, // processed 10% of the time
		},
	})

	mux := asynq.NewServeMux()

	// add handler
	mux.Handle(
		load.KEY,
		load.NewProcessor(
			mongoClient,
			mongoSettings.Database,
			mongoSettings.Collection,
			asynqClient,
		),
	)

	if err := worker.Run(mux); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
