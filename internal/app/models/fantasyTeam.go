package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FantasyTeam struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Date  time.Time          `bson:"date" json:"date"`
	Team  []Player           `bosn:"team" json:"team"`
	Total float32            `bson:"total" json:"total"`
}

type Player struct {
	PlayerCard PlayerCard
	Points     float32
}

func NewTeam() *FantasyTeam {
	return &FantasyTeam{
		ID:   primitive.NewObjectID(),
		Date: time.Now(),
		Team: make([]Player, 5),
	}
}

func (team *FantasyTeam) SetPoints(matches []Match) {
	var total float32
	for _, match := range matches {
		for iPlayer, player := range team.Team {
			if player.PlayerCard.Team == match.Teams[0] || player.PlayerCard.Team == match.Teams[1] {
				for _, v := range match.Points {
					if v.AccountId == player.PlayerCard.AccountId {
						total += v.Total
						team.Team[iPlayer].Points = v.Total
					}
				}
			}
		}
	}

	team.Total = total
}
