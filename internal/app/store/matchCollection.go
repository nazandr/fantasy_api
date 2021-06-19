package store

import (
	"github.com/nazandr/fantasy_api/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MatchCollection struct {
	Store      *Store
	Collection *mongo.Collection
}

func (c *MatchCollection) Create(match *models.Match) error {
	_, err := c.Collection.InsertOne(c.Store.context, match)

	if err != nil {
		return err
	}
	return nil
}

func (c *MatchCollection) GetAll(match *models.Match) ([]models.Match, error) {
	var matches []models.Match
	cursor, err := c.Collection.Find(c.Store.context, bson.D{{}})
	if err != nil {
		return nil, err
	}

	for cursor.Next(c.Store.context) {
		var match models.Match

		if err := cursor.Decode(&match); err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	if cursor.Err() != nil {
		return nil, err
	}

	cursor.Close(c.Store.context)

	return matches, nil
}
