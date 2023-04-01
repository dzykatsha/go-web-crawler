package crawler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/dzykatsha/go-web-crawler/internal/utils"
	"github.com/hibiken/asynq"
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

	// get passed page
	req, err := http.Get(payload.Url)
	if err != nil {
		return err
	}

	// check if we actually work with HTML file
	if !strings.Contains(req.Header.Get("content-type"), "html") {
		return fmt.Errorf("content-type is not html: %s", payload.Url)
	}

	// read response body
	body, err := io.ReadAll(req.Body)
	if err != nil {
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
		return err
	}
	// check if we need to recursively load other urls
	if payload.Depth == 0 {
		return nil
	}
	// parse response body for links
	document, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return err
	}
	urls := utils.FindUrls(document)

	// send recursively other links with lower depth
	for _, nestedUrl := range urls {
		subTask, err := NewLoadURLTask(nestedUrl, payload.Depth-1)
		if err != nil {
			continue
		}

		p.client.Enqueue(subTask)
	}

	return nil
}

func NewLoadURLProcessor(mongoClient *mongo.Client, database string, collection string, asynqClient *asynq.Client) *LoadURLProcessor {
	return &LoadURLProcessor{
		collection: mongoClient.Database(database).Collection(collection),
		client:     asynqClient,
	}
}
