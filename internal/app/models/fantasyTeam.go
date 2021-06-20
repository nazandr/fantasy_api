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
	ID        primitive.ObjectID
	AccountId int
	TeamName  string
	Points    float32
}

func NewTeam() *FantasyTeam {
	return &FantasyTeam{
		ID:   primitive.NewObjectID(),
		Date: time.Now(),
		Team: make([]Player, 5),
	}
}

func (team *FantasyTeam) SetPoints(matches []Match) {
	for _, match := range matches {
		for _, player := range team.Team {
			if player.TeamName == match.Teams[0] || player.TeamName == match.Teams[1] {
				for _, v := range match.Points {
					if v.AccountId == player.AccountId {
						player.Points = v.Total
					}
				}
			}
		}
	}
}
