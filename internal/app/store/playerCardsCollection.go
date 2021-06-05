package store

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlayerCardsCollection struct {
	Store      *Store
	Collection *mongo.Collection
}

type PlayerCard struct {
	Id          primitive.ObjectID `bson:"_id"`
	AccountId   int                `bson:"account_id" json:"account_id"`
	Name        string             `bson:"name" json:"name"`
	FantacyRole int                `bson:"fantasy_role" json:"fantasy_role"`
	Team        string             `bson:"team" json:"team_name"`
	Rarity      int                `json:"rarity"`
}

func NewPlayer() *PlayerCard {
	return &PlayerCard{}
}

func (p *PlayerCardsCollection) GetAll() ([]PlayerCard, error) {
	var playersList []PlayerCard
	cursor, err := p.Collection.Find(p.Store.context, bson.D{{}})
	if err != nil {
		return nil, err
	}

	for cursor.Next(p.Store.context) {
		var player PlayerCard

		if err := cursor.Decode(&player); err != nil {
			return nil, err
		}
		playersList = append(playersList, player)
	}

	if cursor.Err() != nil {
		return nil, err
	}

	cursor.Close(p.Store.context)

	return playersList, nil
}
