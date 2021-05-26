package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                primitive.ObjectID `bson:"_id"`
	Email             string             `bson:"email"`
	EncriptedPassword string             `bson:"encripted_password"`
}
