package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                primitive.ObjectID `bson:"_id" json:"id"`
	Email             string             `bson:"email" json:"email"`
	Password          string             `bosn:"_" json:"password,omitempty"`
	EncryptedPassword string             `bson:"encripted_password" json:"-"`
	FantacyCoins      int                `bson:"fantasy_coins" json:"fantacy_coins"`
	Packs             PacksCount
	CardsCollection   [][]PlayerCard `bson:"card_collection"`
	Session           session
}

type PacksCount struct {
	Common  int `bson:"common" json:"common"`
	Special int `bson:"special" json:"special"`
}

type session struct {
	Refresh_token string    `bson:"refresh_token" json:"resresh_token"`
	Expires_at    time.Time `bson:"expires_at"`
}

func NewUser() *User {
	return &User{
		ID:                primitive.NewObjectID(),
		Email:             "",
		Password:          "",
		EncryptedPassword: "",
		FantacyCoins:      0,
		Packs: PacksCount{
			Common:  5,
			Special: 0,
		},
		CardsCollection: [][]PlayerCard{},
		Session:         session{},
	}
}

func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.By(RequiredIf(u.EncryptedPassword == "")), validation.Length(6, 100)))
}

func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		encript, err := encriptPassword(u.Password)
		if err != nil {
			return err
		}
		u.EncryptedPassword = encript
	}
	return nil
}

func (u *User) Sanitaze() {
	u.Password = ""
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func encriptPassword(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)

	if err != nil {
		return "", err
	}

	return string(b), nil
}
