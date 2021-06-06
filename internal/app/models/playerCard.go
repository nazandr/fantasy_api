package models

import (
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlayerCard struct {
	Id          primitive.ObjectID `bson:"_id"`
	AccountId   int                `bson:"account_id" json:"account_id"`
	Name        string             `bson:"name" json:"name"`
	FantacyRole int                `bson:"fantasy_role" json:"fantasy_role"`
	Team        string             `bson:"team" json:"team_name"`
	Rarity      int                `json:"rarity"`

	Buffs []Buff `json:"buffs"`
}

type Buff struct {
	ID         primitive.ObjectID `bson:"_id"`
	NameOfFild string             `bson:"name_of_fild"`
	Multiplier float64            `bson:"multiplier"`
}

func NewPlayerCard() *PlayerCard {
	return &PlayerCard{
		Id: primitive.NewObjectID(),
	}
}

func RandomWithProbabilitis(items []int, weights []float32) int {
	rand.Seed(time.Now().UnixNano())
	var (
		sumWeight []float32
		sum       float32
	)
	for _, v := range weights {
		sum += v
		sumWeight = append(sumWeight, sum)
	}

	ri := rand.Float32()

	for i, v := range sumWeight {
		if v >= ri {
			return items[i]
		}
	}
	return 0
}
