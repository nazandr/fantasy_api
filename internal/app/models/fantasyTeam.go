package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FantasyTeam struct {
	ID   primitive.ObjectID `bson:"_id" json:"id"`
	Date time.Time          `bson:"date" json:"date"`
	Team []Player           `bosn:"team" json:"team"`
}

type Player struct {
	PlayerID primitive.ObjectID
	Points   float32
}

func NewTeam() *FantasyTeam {
	return &FantasyTeam{
		ID:   primitive.NewObjectID(),
		Date: time.Now(),
		Team: make([]Player, 5),
	}
}
