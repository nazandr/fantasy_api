package models

import (
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PlayerCard struct {
	Id          primitive.ObjectID `bson:"_id" json:"id"`
	AccountId   int                `bson:"account_id" json:"account_id"`
	Name        string             `bson:"name" json:"name"`
	FantasyRole int                `bson:"fantasy_role" json:"fantasy_role"`
	Team        string             `bson:"team" json:"team_name"`
	Rarity      int                `json:"rarity"`

	Buffs []Buff `json:"buffs"`
}

type Buff struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	NameOfFild    string             `bson:"name_of_fild" json:"name_of_fild"`
	DisplayedName string             `bson:"displayed_name" json:"displayed_name"`
	Multiplier    int                `bson:"multiplier" json:"multiplier"`
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

func SetBuffs(n int) []Buff {
	rand.Seed(time.Now().UnixNano())
	filds := []string{"Kills", "Assists", "LastHits", "Gpm", "TowerKills", "RoshanKils",
		"Participation", "Observers", "CampStacked", "Runs", "FirsBlood", "Stuns"}
	nameFilds := []string{"Убийства", "Помощь", "Добито", "GPM", "Башни", "Рошаны",
		"Участие в убийствах", "Варды", "Стаки", "Руны", "Первая кровь", "Станы"}
	mult := []int{5, 10, 15, 20, 25}
	buffs := []Buff{}
	for i := 0; i < n; i++ {
		buff := Buff{}
		buff.ID = primitive.NewObjectID()
		r := rand.Intn(len(filds))
		buff.NameOfFild = filds[r]
		buff.DisplayedName = nameFilds[r]
		buff.Multiplier = mult[rand.Intn(len(mult))]
		buffs = append(buffs, buff)
		filds = removeFild(filds, r)
		nameFilds = removeFild(nameFilds, r)
	}

	return buffs
}

func removeFild(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
