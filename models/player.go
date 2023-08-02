package models

import "gopkg.in/mgo.v2/bson"

type Player struct {
	Id      bson.ObjectId `json:"id" bson:"_id"`
	Name    string        `json:"name" bson:"name"`
	Country string        `json:"country" bson:"country"`
	Score   int           `json:"score" bson:"score"`
}

