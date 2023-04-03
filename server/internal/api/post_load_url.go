package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dzykatsha/go-web-crawler/internal/crawler/load"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type PostLoadURLHandler struct {
	asynqClient *asynq.Client
}

type PostLoadURLData struct {
	Url   string `json:"url"`
	Depth int    `json:"depth"`
}

func (handler PostLoadURLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// verify method
	if r.Method != "POST" {
		log.Error().Msgf("method not allowed: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method not allowed"))
		return
	}

	// read body
	rawRequestBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msgf("failed to read body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("400 - Failed to read body\n%v", err)))
		return
	}

	var requestBody PostLoadURLData
	if err := json.Unmarshal(rawRequestBody, &requestBody); err != nil {
		log.Error().Err(err).Msgf("failed to parse body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("400 - Failed to parse body\n%v", err)))
		return
	}

	// send task
	task, err := load.NewTask(requestBody.Url, requestBody.Depth)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create task")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("500 - failed to create task: %v", err)))
		return
	}

	taskInfo, err := handler.asynqClient.Enqueue(task, asynq.MaxRetry(-1))
	if err != nil {
		log.Error().Err(err).Msgf("failed send load url task")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("500 - failed to send load url task: %v", err)))
		return
	}

	// respond
	log.Info().Msg("Successfully send load url task")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(taskInfo.ID))
}

func NewPostLoadURLHandler(client *asynq.Client) PostLoadURLHandler {
	return PostLoadURLHandler{
		asynqClient: client,
	}
}
