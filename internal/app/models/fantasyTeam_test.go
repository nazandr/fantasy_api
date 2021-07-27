package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFantasyTeam_SetPoints(t *testing.T) {
	type fields struct {
		ID         primitive.ObjectID
		Date       time.Time
		Team       []Player
		Calculated bool
		Total      float32
	}
	type args struct {
		serires []Series
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		value  float32
	}{
		{
			name: "valid",
			fields: fields{
				Team: []Player{
					{
						PlayerCard: PlayerCard{
							AccountId: 1,
							Team:      "A",
							Buffs: []Buff{
								{
									NameOfFild: "Kills",
									Multiplier: 10,
								},
							},
						},
						Points: 0,
					},
				},
			},
			args: args{
				serires: []Series{
					{
						Teams: []string{"A", "B"},
						Matches: []Match{
							{
								Teams: []string{"A", "B"},
								Points: []Points{
									{
										AccountId: 1,
										Kills:     10,
										Total:     8,
									},
								},
							},
							{
								Teams: []string{"A", "B"},
								Points: []Points{
									{
										AccountId: 1,
										Kills:     10,
										Total:     10,
									},
								},
							},
							{
								Teams: []string{"A", "B"},
								Points: []Points{
									{
										AccountId: 1,
										Kills:     10,
										Total:     9,
									},
								},
							},
						},
					},
				},
			},
			value: 21,
		},
		{
			name: "one match",
			fields: fields{
				Team: []Player{
					{
						PlayerCard: PlayerCard{
							AccountId: 1,
							Team:      "A",
							Buffs: []Buff{
								{
									NameOfFild: "Kills",
									Multiplier: 10,
								},
							},
						},
						Points: 0,
					},
				},
			},
			args: args{
				serires: []Series{
					{
						Teams: []string{"A", "B"},
						Matches: []Match{
							{
								Teams: []string{"A", "B"},
								Points: []Points{
									{
										AccountId: 1,
										Kills:     10,
										Total:     8,
									},
								},
							},
						},
					},
				},
			},
			value: 9,
		},
		{
			name: "noone",
			fields: fields{
				Team: []Player{
					{
						PlayerCard: PlayerCard{
							AccountId: 1,
							Team:      "A",
							Buffs: []Buff{
								{
									NameOfFild: "Kills",
									Multiplier: 10,
								},
							},
						},
						Points: 0,
					},
				},
			},
			args: args{
				serires: []Series{
					{
						Teams: []string{"A", "B"},
						Matches: []Match{
							{
								Teams: []string{"A", "B"},
								Points: []Points{
									{
										AccountId: 2,
										Kills:     10,
										Total:     8,
									},
								},
							},
						},
					},
				},
			},
			value: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			team := &FantasyTeam{
				ID:         tt.fields.ID,
				Date:       tt.fields.Date,
				Team:       tt.fields.Team,
				Calculated: tt.fields.Calculated,
				Total:      tt.fields.Total,
			}
			team.SetPoints(tt.args.serires)
			assert.Equal(t, tt.value, team.Total)
		})
	}
}
