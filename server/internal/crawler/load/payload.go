package load

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const KEY = "crawler:loadUrl"

type Payload struct {
	Url   string
	Depth int
}

func NewTask(url string, depth int) (*asynq.Task, error) {
	payload, err := json.Marshal(Payload{Url: url, Depth: depth})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(KEY, payload), nil
}
