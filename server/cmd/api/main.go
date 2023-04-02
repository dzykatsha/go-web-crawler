package main

import (
	"context"
	"net/http"
	"time"

	"github.com/dzykatsha/go-web-crawler/internal/api"
	"github.com/dzykatsha/go-web-crawler/internal/settings"
	"github.com/hibiken/asynq"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// setup settings
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)
	redisSettings := settings.NewRedisSettingsFromEnv()
	apiSettings := settings.NewAPISettingsFromEnv()
	mongoSettings := settings.NewMongoSettingsFromEnv()

	// setup asynq
	redisConnection := asynq.RedisClientOpt{
		Addr:     redisSettings.Address(),
		Password: redisSettings.Password,
	}

	asynqClient := asynq.NewClient(redisConnection)
	defer asynqClient.Close()

	// setup mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoSettings.ConnectionURL()))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to mongodb")
	}

	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			log.Panic().Err(err)
		}
	}()
	collection := mongoClient.Database(mongoSettings.Database).Collection(mongoSettings.Collection)

	// setup http
	mux := http.NewServeMux()

	mux.Handle("/load", api.NewPostLoadURLHandler(asynqClient))
	mux.Handle("/statistics", api.NewGetStatisticsHandler(collection))
	mux.Handle("/page", api.NewGetPageHandler(collection))

	handler := cors.Default().Handler(mux)
	if err := http.ListenAndServe(apiSettings.Address(), handler); err != nil {
		log.Panic().Err(err)
	}
}
