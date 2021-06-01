package store

import (
	"errors"
	"time"

	"github.com/nazandr/fantasy_api/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserCollection struct {
	Store      *Store
	Collection *mongo.Collection
}

var (
	ErrUserAllreadyExist = errors.New("User allready exist")
)

func (c *UserCollection) Create(u *models.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	if _, err := c.FindByEmail(u.Email); err != mongo.ErrNoDocuments {
		return ErrUserAllreadyExist
	}

	res, err := c.Collection.InsertOne(c.Store.context, bson.D{
		{Key: "email", Value: u.Email},
		{Key: "encripted_password", Value: u.EncryptedPassword},
	})

	if err != nil {
		return err
	}

	u.ID = res.InsertedID.(primitive.ObjectID)
	return nil
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

func (c *UserCollection) UpdateRefreshToken(id primitive.ObjectID, rt string, exp int) error {
	_, err := c.Collection.UpdateByID(
		c.Store.context,
		id,
		bson.D{
			{"$set", bson.D{{"session", bson.D{
				{"refresh_token", rt},
				{"expires_at", time.Now().Add(time.Minute * time.Duration(exp))},
			}}},
			}},
	)
	if err != nil {
		return err
	}

	return nil
}
