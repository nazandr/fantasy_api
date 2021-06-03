package models

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUser(t *testing.T) *User {
	return &User{
		ID:       primitive.NewObjectID(),
		Email:    "user@example.com",
		Password: "password",
		Packs: Packs{
			Common:  0,
			Special: 0,
		},
		Session: session{
			Refresh_token: "refresh token",
			Expires_at:    time.Now(),
		},
	}
}
