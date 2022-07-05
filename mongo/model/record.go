package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type RecordCache struct {
	Id      primitive.ObjectID `json:"id" bson:"_id"`
	Room    string             `bson:"room"`
	Records []Record           `bson:"records"`
}

type Record struct {
	Display   string `bson:"display"`
	StartTime int64  `bson:"startTime"`
	End       bool   `bson:"end"`
	Url       string `bson:"url"`
}
