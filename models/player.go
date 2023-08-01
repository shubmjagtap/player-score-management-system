package models

import "gopkg.in/mgo.v2/bson"

type Player struct {
	Id      bson.ObjectId `json:"id" bson:"_id"`
	Name    string        `json:"name" bson:"name"`
	Country string        `json:"country" bson:"country"`
	Score   int           `json:"score" bson:"score"`
}

// All the attributes are mandatory.
// ID uniquely identifies each player.
// Name should have a cap of 15 characters.
// Country code will be a two letter code representing the country (For e.g., IN, US).
// Score will be an integer representing player score
