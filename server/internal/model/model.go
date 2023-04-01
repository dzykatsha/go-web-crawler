package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URLDocument struct {
	ID        primitive.ObjectID `bson:"_id" json:"uid"`
	CreatedAt int64              `bson:"createdAt" json:"createdAt"`
	HTML      string             `bson:"html" json:"html"`
	Url       string             `bson:"url" json:"url"`
}
