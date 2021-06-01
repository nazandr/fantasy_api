package api_server

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	AcssesToken  string `json:"acsses_token"`
	RefreshToken string `json:"refresh_token"`
}

var ()

func NewToken() *Token {
	return &Token{
		AcssesToken:  "",
		RefreshToken: "",
	}
}

func (t *Token) NewJWT(id primitive.ObjectID, conf *Config) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(conf.AcssesTokenExp)).Unix(),
		Subject:   id.String(),
	})
	at, err := token.SignedString([]byte(conf.SignatureKey))
	if err != nil {
		return err
	}
	t.AcssesToken = at

	return nil
}

func (t *Token) NewRefreshToken() error {
	b := make([]byte, 32)
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	_, err := r.Read(b)
	if err != nil {
		return err
	}

	t.RefreshToken = string(b)
	return nil
}

func (t *Token) Auth(id primitive.ObjectID, conf *Config) error {
	if err := t.NewJWT(id, conf); err != nil {
		return err
	}

	if err := t.NewRefreshToken(); err != nil {
		return err
	}

	return nil
}

func (t *Token) ParseJWT(conf *Config) (string, error) {
	token, err := jwt.Parse(t.AcssesToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(conf.SignatureKey), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["sub"].(string), nil
}
