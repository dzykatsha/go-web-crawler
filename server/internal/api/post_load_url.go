package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dzykatsha/go-web-crawler/internal/crawler"
	"github.com/hibiken/asynq"
)

type PostLoadURLHandler struct {
	asynqClient *asynq.Client
}

type PostLoadURLData struct {
	Url   string `json:"url"`
	Depth int    `json:"depth"`
}

func (h PostLoadURLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method not allowed"))
		return
	}

	var b PostLoadURLData
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("400 - Failed to read body\n%v", err)))
		return
	}

	if err := json.Unmarshal(raw, &b); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("400 - Failed to read body\n%v", err)))
		return
	}

	if b.Url == "" || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - missing url or depth query parameter"))
		return
	}

	t, err := crawler.NewLoadURLTask(b.Url, b.Depth)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("500 - failed to parse params: %v", err)))
		return
	}

	i, err := h.asynqClient.Enqueue(t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("500 - failed to send load url task: %v", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(i.ID))
}

func NewPostLoadURLHandler(client *asynq.Client) PostLoadURLHandler {
	return PostLoadURLHandler{
		asynqClient: client,
	}
}
