package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FantacyTeam struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Date time.Time          `bson:"date" json:"date"`
	Team []PlayerCard       `bosn:"team" json:"team"`
}

func NewTeam() *FantacyTeam {
	return &FantacyTeam{
		ID:   primitive.NewObjectID(),
		Date: time.Now(),
		Team: make([]PlayerCard, 5),
	}
}
