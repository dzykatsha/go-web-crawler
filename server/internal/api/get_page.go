package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dzykatsha/go-web-crawler/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetPageHandler struct {
	collection *mongo.Collection
}

func (handler GetPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// verify method
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method not allowed"))
		return
	}

	// get uid from query
	query := r.URL.Query()
	uid := query.Get("uid")
	if uid == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - No uid in query params"))
		return
	}

	// get data by uid
	objectID, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("400 - Wrong uid\n%v", err)))
		return
	}

	singleResult := handler.collection.FindOne(context.TODO(), bson.M{"_id": objectID})
	err = singleResult.Err()
	if err != nil {
		switch err.Error() {
		case mongo.ErrNoDocuments.Error():
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("404 - Not found document: %s", uid)))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("500 - Failed to fetch data\n%v", err)))
			return
		}
	}

	var url model.URLDocument
	err = singleResult.Decode(&url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("500 - Failed to fetch data\n%v", err)))
		return
	}

	// respond
	responseBody, err := json.Marshal(url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("500 - Failed to fetch data\n%v", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

func NewGetPageHandler(collection *mongo.Collection) GetPageHandler {
	return GetPageHandler{
		collection: collection,
	}
}
