package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dzykatsha/go-web-crawler/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GetStatisticsHandler struct {
	collection *mongo.Collection
}

type GetStatisticsData struct {
	Total int64               `json:"total"`
	Urls  []model.URLDocument `json:"urls"`
}

func (h GetStatisticsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method not allowed"))
		return
	}

	q := r.URL.Query()
	p, err := strconv.Atoi(q.Get("page"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Failed to parse int from page query parameter"))
		return
	}
	o := options.Find().
		SetSort(bson.M{"createdAt": 1}).
		SetSkip((int64)(p-1) * 10).
		SetLimit(10)
	ctx := context.TODO()
	c, err := h.collection.Find(ctx, bson.D{}, o)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("500 - Failed to fetch data: %v", err)))
		return
	}

	var result []model.URLDocument
	if err := c.All(ctx, &result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("500 - Failed to fetch data: %v", err)))
		return
	}
	total, err := h.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("500 - Failed to fetch total: %v", err)))
		return
	}

	resBody, err := json.Marshal(GetStatisticsData{Total: total, Urls: result})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("500 - Failed to marshal data: %v", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}

func NewGetStatisticsHandler(collection *mongo.Collection) GetStatisticsHandler {
	return GetStatisticsHandler{
		collection: collection,
	}
}
