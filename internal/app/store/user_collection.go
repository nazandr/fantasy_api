package store

import (
	"github.com/nazandr/fantasy_api/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserCollection struct {
	Store      *Store
	Collection *mongo.Collection
}

func (c *UserCollection) Create(u *models.User) (*models.User, error) {
	u.ID = primitive.NewObjectID()
	if _, err := c.Collection.InsertOne(c.Store.context, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (c *UserCollection) FindByEmail(email string) (*models.User, error) {
	u := &models.User{}

	filter := bson.D{primitive.E{
		Key:   "email",
		Value: email,
	}}

	if err := c.Collection.FindOne(c.Store.context, filter).Decode(&u); err != nil {
		return nil, err
	}

	return u, nil
}
