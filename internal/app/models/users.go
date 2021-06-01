package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                primitive.ObjectID `bson:"_id" json:"id"`
	Email             string             `bson:"email" json:"email"`
	Password          string             `json:"password,omitempty"`
	EncryptedPassword string             `bson:"encripted_password" json:"-"`
	RefreshToken      string             `bson:"refresh_token" json:"resresh_token"`
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
