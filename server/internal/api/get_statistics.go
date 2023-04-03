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

func (handler GetStatisticsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// verify method
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method not allowed"))
		return
	}

	// read page from query
	query := r.URL.Query()
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Failed to parse int from page query parameter"))
		return
	}

	// read documents from mongo
	mongoOptions := options.Find().
		SetSort(bson.M{"createdAt": 1}).
		SetSkip((int64)(page-1) * 10).
		SetLimit(10)
	ctx := context.TODO()
	mongoCursor, err := handler.collection.Find(ctx, bson.D{}, mongoOptions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("500 - Failed to fetch data: %v", err)))
		return
	}

	var result []model.URLDocument
	if err := mongoCursor.All(ctx, &result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("500 - Failed to fetch data: %v", err)))
		return
	}

	// read total count from mongo
	total, err := handler.collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("500 - Failed to fetch total: %v", err)))
		return
	}

	// respond
	responseBody, err := json.Marshal(GetStatisticsData{Total: total, Urls: result})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("500 - Failed to marshal data: %v", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

func NewGetStatisticsHandler(collection *mongo.Collection) GetStatisticsHandler {
	return GetStatisticsHandler{
		collection: collection,
	}
}
