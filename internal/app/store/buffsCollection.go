package store

import (
	"github.com/nazandr/fantasy_api/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BuffsCollection struct {
	Store      *Store
	Collection *mongo.Collection
}

func NewBuff() *models.Buff {
	return &models.Buff{
		ID: primitive.NewObjectID(),
	}
}

func (b *BuffsCollection) GetAll() ([]models.Buff, error) {
	var buffsList []models.Buff
	cursor, err := b.Collection.Find(b.Store.context, bson.D{{}})
	if err != nil {
		return nil, err
	}

	for cursor.Next(b.Store.context) {
		var buff models.Buff

		if err := cursor.Decode(&buff); err != nil {
			return nil, err
		}
		buffsList = append(buffsList, buff)
	}

	if cursor.Err() != nil {
		return nil, err
	}

	cursor.Close(b.Store.context)

	return buffsList, nil
}
