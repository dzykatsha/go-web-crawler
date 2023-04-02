package crawler

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

type LoadURLProcessor struct {
	collection *mongo.Collection
	client     *asynq.Client
}

func (p LoadURLProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	// read payload
	var payload LoadURLPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	log.Info().
		Str("url", payload.Url).
		Int("depth", payload.Depth).
		Msg("received new url to parse")

	baseUrl, err := url.Parse(payload.Url)
	if err != nil {
		log.Error().
			Str("url", payload.Url).
			Int("depth", payload.Depth).
			Err(err).
			Msg("wrong url")

		return err
	}

	// get passed page
	req, err := http.Get(payload.Url)
	if err != nil {
		log.Error().
			Str("url", payload.Url).
			Int("depth", payload.Depth).
			Err(err).
			Msg("failed to get page")

		return err
	}

	// check if we actually work with HTML file
	if !strings.Contains(req.Header.Get("content-type"), "html") {
		log.Error().
			Str("url", payload.Url).
			Int("depth", payload.Depth).
			Msg("not html")

		return fmt.Errorf("content-type is not html: %s", payload.Url)
	}

	// read response body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Error().
			Str("url", payload.Url).
			Int("depth", payload.Depth).
			Err(err).
			Msg("failed to read response body")

		return err
	}

	// save to mongodb
	data := bson.M{
		"createdAt": time.Now().Unix(),
		"html":      string(body),
		"url":       payload.Url,
	}

	_, err = p.collection.InsertOne(ctx, data)
	if err != nil {
		log.Error().
			Str("url", payload.Url).
			Int("depth", payload.Depth).
			Err(err).
			Msg("failed to inset document")

		return err
	}
	// check if we need to recursively load other urls
	if payload.Depth == 0 {
		log.Info().
			Str("url", payload.Url).
			Int("depth", payload.Depth).
			Msg("finished processing leaf url")

		return nil
	}
	// parse response body for links

	urls := []string{}

	for _, value := range xurls.Relaxed.FindAllString(string(body), -1) {
		u, err := url.Parse(value)
		if err != nil {
			log.Error().
				Str("url", payload.Url).
				Int("depth", payload.Depth).
				Err(err).
				Msg("failed to parse next url")

			continue
		}

		if u.Scheme == "" || u.Host == "" {
			u.Scheme = baseUrl.Scheme
			u.Host = baseUrl.Host
		}

		if u.Host != baseUrl.Host {
			log.Error().
				Str("url", payload.Url).
				Int("depth", payload.Depth).
				Msg("host mismatch")

			continue
		}

		urls = append(urls, u.String())
	}

	log.Info().
		Int("count", len(urls)).
		Msg("invoking urls")

	// send recursively other links with lower depth
	for _, nestedUrl := range urls {
		subTask, err := NewLoadURLTask(nestedUrl, payload.Depth-1)
		if err != nil {
			continue
		}

		p.client.Enqueue(subTask, asynq.MaxRetry(-1))
	}
	log.Info().
		Str("url", payload.Url).
		Int("depth", payload.Depth).
		Msg("finished processing node url")

	return nil
}

func NewLoadURLProcessor(mongoClient *mongo.Client, database string, collection string, asynqClient *asynq.Client) *LoadURLProcessor {
	return &LoadURLProcessor{
		collection: mongoClient.Database(database).Collection(collection),
		client:     asynqClient,
	}
}
