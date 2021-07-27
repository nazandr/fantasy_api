package models

import (
	"math"
	"reflect"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FantasyTeam struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Date       time.Time          `bson:"date" json:"date"`
	Team       []Player           `bosn:"team" json:"team"`
	Calculated bool               `bson:"calculated" json:"-"`
	Total      float32            `bson:"total" json:"total"`
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

func (team *FantasyTeam) SetPoints(serires []Series) {
	for _, ser := range serires {
		for iPlayer, player := range team.Team {
			if ser.Teams[0] == player.PlayerCard.Team || ser.Teams[1] == player.PlayerCard.Team {
				selected := []Player{}

				for _, match := range ser.Matches {
					for _, points := range match.Points {
						if points.AccountId == player.PlayerCard.AccountId {
							p := points.Total
							if len(player.PlayerCard.Buffs) > 0 {
								for _, buff := range player.PlayerCard.Buffs {
									i := reflect.ValueOf(points).FieldByName(buff.NameOfFild).Float()
									p += float32(i * (float64(buff.Multiplier) / 100))
								}
							}
							selected = append(selected, Player{
								PlayerCard: player.PlayerCard,
								Points:     p,
							})
						}
					}
				}

				switch len(selected) {
				case 0:
					continue
				case 1:
					selected = append(selected, Player{})
				}

				sort.SliceStable(selected, func(i, j int) bool {
					return selected[i].Points > selected[j].Points
				})

				for i := 0; i < 2; i++ {
					team.Team[iPlayer].Points += selected[i].Points
					team.Total += selected[i].Points
				}
			}
		}
	}

	for i := 0; i < len(team.Team); i++ {
		team.Team[i].Points = float32(math.Round(float64(team.Team[i].Points*100)) / 100)
	}

	team.Total = float32(math.Round(float64(team.Total*100)) / 100)
}
