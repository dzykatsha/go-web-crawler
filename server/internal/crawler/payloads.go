package crawler

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const LoadURL = "crawler:loadUrl"

type LoadURLPayload struct {
	Url   string
	Depth int
}

func NewLoadURLTask(url string, depth int) (*asynq.Task, error) {
	payload, err := json.Marshal(LoadURLPayload{Url: url, Depth: depth})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(LoadURL, payload), nil
}
