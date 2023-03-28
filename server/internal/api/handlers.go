package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dzykatsha/go-web-crawler/internal/crawler"
	"github.com/hibiken/asynq"
)

type PostLoadURLHandler struct {
	asynqClient *asynq.Client
}

func (h PostLoadURLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method not allowed"))
		return
	}

	q := r.URL.Query()
	url := q.Get("url")
	depth, err := strconv.Atoi(q.Get("depth"))

	if url == "" || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - missing url or depth query parameter"))
		return
	}

	t, err := crawler.NewLoadURLTask(url, depth)
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
