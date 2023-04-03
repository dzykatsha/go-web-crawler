package load

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/mvdan/xurls"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Processor struct {
	collection *mongo.Collection
	client     *asynq.Client
}

func (p Processor) ProcessTask(ctx context.Context, task *asynq.Task) error {
	// read payload
	var payload Payload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	taskId := task.ResultWriter().TaskID()

	log.Info().
		Str("task", taskId).
		Str("url", payload.Url).
		Int("depth", payload.Depth).
		Msg("received new url to parse")

	baseUrl, err := url.Parse(payload.Url)
	if err != nil {
		log.Error().
			Str("task", taskId).
			Str("url", payload.Url).
			Int("depth", payload.Depth).
			Err(err).
			Msg("wrong url")

		return err
	}

	// get passed page
	response, err := http.Get(payload.Url)
	if err != nil {
		log.Error().
			Str("task", taskId).
			Str("url", payload.Url).
			Int("depth", payload.Depth).
			Err(err).
			Msg("failed to get page")

		return err
	}

	// check if we actually work with HTML file
	if !strings.Contains(response.Header.Get("content-type"), "html") {
		log.Error().
			Str("task", taskId).
			Str("url", payload.Url).
			Int("depth", payload.Depth).
			Msg("not html")

		return fmt.Errorf("content-type is not html: %s", payload.Url)
	}

	// read response responseBody
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error().
			Str("task", taskId).
			Str("url", payload.Url).
			Int("depth", payload.Depth).
			Err(err).
			Msg("failed to read response body")

		return err
	}

	// save to mongodb
	data := bson.M{
		"createdAt": time.Now().Unix(),
		"html":      string(responseBody),
		"url":       payload.Url,
	}

	_, err = p.collection.InsertOne(ctx, data)
	if err != nil {
		log.Error().
			Str("task", taskId).
			Str("url", payload.Url).
			Int("depth", payload.Depth).
			Err(err).
			Msg("failed to insert document")

		return err
	}
	// check if we need to recursively load other urls
	if payload.Depth == 0 {
		log.Info().
			Str("task", taskId).
			Str("url", payload.Url).
			Int("depth", payload.Depth).
			Msg("finished processing leaf url")

		return nil
	}
	// parse response body for links

	urls := []string{}

	for _, value := range xurls.Relaxed.FindAllString(string(responseBody), -1) {
		url, err := ParseUrl(value, *baseUrl)
		if err != nil {
			log.Error().
				Str("task", taskId).
				Str("url", payload.Url).
				Int("depth", payload.Depth).
				Err(err).
				Msgf("failed to parse url: %s", url)

			continue
		}

		urls = append(urls, url)
	}

	log.Info().
		Str("task", taskId).
		Str("url", payload.Url).
		Int("depth", payload.Depth).
		Int("count", len(urls)).
		Msg("invoking urls")

	// send recursively other links with lower depth
	for _, url := range urls {
		subTask, err := NewTask(url, payload.Depth-1)
		if err != nil {
			continue
		}

		p.client.Enqueue(subTask, asynq.MaxRetry(-1))
	}
	log.Info().
		Str("task", taskId).
		Str("url", payload.Url).
		Int("depth", payload.Depth).
		Msg("finished processing node url")

	return nil
}

func NewProcessor(mongoClient *mongo.Client, database string, collection string, asynqClient *asynq.Client) *Processor {
	return &Processor{
		collection: mongoClient.Database(database).Collection(collection),
		client:     asynqClient,
	}
}
